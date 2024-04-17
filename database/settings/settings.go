// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/util"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

const (
	TableSettings = "settings"
)

var (
	// ErrEmptyCloneImage defines the error type when a
	// Settings type has an empty CloneImage field provided.
	ErrEmptyCloneImage = errors.New("empty settings clone image provided")
	// ErrEmptyQueueRoutes defines the error type when a
	// Settings type has an empty QueueRoutes field provided.
	ErrEmptyQueueRoutes = errors.New("empty settings queue routes provided")
)

// todo: comments Build->Settings
type (
	// config represents the settings required to create the engine that implements the BuildInterface interface.
	config struct {
		// specifies to skip creating tables and indexes for the Build engine
		SkipCreation bool
	}

	// engine represents the build functionality that implements the BuildInterface interface.
	engine struct {
		// engine configuration settings used in build functions
		config *config

		ctx context.Context

		// gorm.io/gorm database client used in build functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in build functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}

	// Settings is the database representation of platform settings.
	Settings struct {
		ID                sql.NullInt64  `sql:"id"`
		CloneImage        sql.NullString `sql:"clone_image"`
		TemplateDepth     sql.NullInt64  `sql:"template_depth"`
		StarlarkExecLimit sql.NullInt64  `sql:"starlark_exec_limit"`
		RepoAllowlist     pq.StringArray `sql:"repo_allowlist" gorm:"type:varchar(1000)"`
		QueueRoutes       pq.StringArray `sql:"queue_routes" gorm:"type:varchar(1000)"`
	}
)

// New creates and returns a Vela service for integrating with builds in the database.
//
//nolint:revive // ignore returning unexported engine
func New(opts ...EngineOpt) (*engine, error) {
	// create new Build engine
	e := new(engine)

	// create new fields
	e.client = new(gorm.DB)
	e.config = new(config)
	e.logger = new(logrus.Entry)

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(e)
		if err != nil {
			return nil, err
		}
	}

	// check if we should skip creating database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of settings table and indexes in the database")

		return e, nil
	}

	// create the settings table
	err := e.CreateSettingsTable(e.ctx, e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", TableSettings, err)
	}

	// todo: need indexes?

	return e, nil
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
//
// When a field within the Settings type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (s *Settings) Nullify() *Settings {
	if s == nil {
		return nil
	}

	// check if the ID field should be false
	if s.ID.Int64 == 0 {
		s.ID.Valid = false
	}

	// check if the CloneImage field should be false
	if len(s.CloneImage.String) == 0 {
		s.CloneImage.Valid = false
	}

	return s
}

// ToAPI converts the Settings type
// to an API Settings type.
func (s *Settings) ToAPI() *api.Settings {
	settings := new(api.Settings)

	settings.SetID(s.ID.Int64)
	settings.SetCloneImage(s.CloneImage.String)
	settings.SetTemplateDepth(int(s.TemplateDepth.Int64))
	settings.SetStarlarkExecLimit(uint64(s.StarlarkExecLimit.Int64))
	settings.SetCloneImage(s.CloneImage.String)
	settings.SetQueueRoutes(s.QueueRoutes)

	return settings
}

// Validate verifies the necessary fields for
// the Settings type are populated correctly.
func (s *Settings) Validate() error {
	// verify the CloneImage field is populated
	if len(s.CloneImage.String) == 0 {
		return ErrEmptyCloneImage
	}

	// ensure that all Settings string fields
	// that can be returned as JSON are sanitized
	// to avoid unsafe HTML content
	s.CloneImage = sql.NullString{String: util.Sanitize(s.CloneImage.String), Valid: s.CloneImage.Valid}

	// ensure that all QueueRoutes are sanitized
	// to avoid unsafe HTML content
	for i, v := range s.QueueRoutes {
		s.QueueRoutes[i] = util.Sanitize(v)
	}

	return nil
}

// FromAPI converts the API Settings type
// to a database Settings type.
func FromAPI(s *api.Settings) *Settings {
	settings := &Settings{
		ID:                sql.NullInt64{Int64: s.GetID(), Valid: true},
		CloneImage:        sql.NullString{String: s.GetCloneImage(), Valid: true},
		TemplateDepth:     sql.NullInt64{Int64: int64(s.GetTemplateDepth()), Valid: true},
		StarlarkExecLimit: sql.NullInt64{Int64: int64(s.GetStarlarkExecLimit()), Valid: true},
		QueueRoutes:       pq.StringArray(s.GetQueueRoutes()),
	}

	return settings.Nullify()
}
