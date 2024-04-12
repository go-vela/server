// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/util"
)

var (
	// ErrEmptyDashName defines the error type when a
	// User type has an empty Name field provided.
	ErrEmptyDashName = errors.New("empty dashboard name provided")

	// ErrExceededAdminLimit defines the error type when a
	// User type has Admins field provided that exceeds the database limit.
	ErrExceededAdminLimit = errors.New("exceeded admins limit")
)

type (
	// config represents the settings required to create the engine that implements the BuildExecutableService interface.
	config struct {
		// specifies to skip creating tables and indexes for the BuildExecutable engine
		SkipCreation bool
		// specifies the driver for proper popping query
		Driver string
	}

	// engine represents the build executable functionality that implements the BuildExecutableService interface.
	engine struct {
		// engine configuration settings used in build executable functions
		config *config

		ctx context.Context

		// gorm.io/gorm database client used in build executable functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in build executable functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}

	// Dashboard is the database representation of a user.
	Dashboard struct {
		ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
		Name      sql.NullString `sql:"name"`
		CreatedAt sql.NullInt64  `sql:"created_at"`
		CreatedBy sql.NullString `sql:"created_by"`
		UpdatedAt sql.NullInt64  `sql:"updated_at"`
		UpdatedBy sql.NullString `sql:"updated_by"`
		Admins    pq.StringArray `sql:"admins" gorm:"type:varchar(5000)"`
		Repos     DashReposJSON
	}

	DashReposJSON []*api.DashboardRepo
)

// New creates and returns a Vela service for integrating with build executables in the database.
//
//nolint:revive // ignore returning unexported engine
func New(opts ...EngineOpt) (*engine, error) {
	// create new BuildExecutable engine
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

	// check if we should skip creating build executable database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of build executables table and indexes in the database")

		return e, nil
	}

	// create the build executables table
	err := e.CreateDashboardTable(e.ctx, e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableDashboard, err)
	}

	return e, nil
}

// Value - Implementation of valuer for database/sql for DashReposJSON.
func (r DashReposJSON) Value() (driver.Value, error) {
	valueString, err := json.Marshal(r)
	return string(valueString), err
}

// Scan - Implement the database/sql scanner interface for DashReposJSON.
func (r *DashReposJSON) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &r)
	case string:
		return json.Unmarshal([]byte(v), &r)
	default:
		return fmt.Errorf("wrong type for repos: %T", v)
	}
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
//
// When a field within the Dashboard type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (d *Dashboard) Nullify() *Dashboard {
	if d == nil {
		return nil
	}

	// check if the Name field should be false
	if len(d.Name.String) == 0 {
		d.Name.Valid = false
	}

	// check if the CreatedAt field should be false
	if d.CreatedAt.Int64 == 0 {
		d.CreatedAt.Valid = false
	}

	// check if the CreatedBy field should be false
	if len(d.CreatedBy.String) == 0 {
		d.CreatedBy.Valid = false
	}

	// check if the UpdatedAt field should be false
	if d.UpdatedAt.Int64 == 0 {
		d.UpdatedAt.Valid = false
	}

	// check if the UpdatedBy field should be false
	if len(d.UpdatedBy.String) == 0 {
		d.UpdatedBy.Valid = false
	}

	return d
}

// ToLibrary converts the Dashboard type
// to a library Dashboard type.
func (d *Dashboard) ToAPI() *api.Dashboard {
	dashboard := new(api.Dashboard)

	dashboard.SetID(d.ID.String())
	dashboard.SetName(d.Name.String)
	dashboard.SetAdmins(d.Admins)
	dashboard.SetCreatedAt(d.CreatedAt.Int64)
	dashboard.SetCreatedBy(d.CreatedBy.String)
	dashboard.SetUpdatedAt(d.UpdatedAt.Int64)
	dashboard.SetUpdatedBy(d.UpdatedBy.String)
	dashboard.SetRepos(d.Repos)

	return dashboard
}

// Validate verifies the necessary fields for
// the Dashboard type are populated correctly.
func (d *Dashboard) Validate() error {
	// verify the Name field is populated
	if len(d.Name.String) == 0 {
		return ErrEmptyDashName
	}

	// calculate total size of favorites
	total := 0
	for _, f := range d.Admins {
		total += len(f)
	}

	// verify the Favorites field is within the database constraints
	// len is to factor in number of comma separators included in the database field,
	// removing 1 due to the last item not having an appended comma
	if (total + len(d.Admins) - 1) > constants.DashboardAdminMaxSize {
		return ErrExceededAdminLimit
	}

	// ensure that all Dashboard string fields
	// that can be returned as JSON are sanitized
	// to avoid unsafe HTML content
	d.Name = sql.NullString{String: util.Sanitize(d.Name.String), Valid: d.Name.Valid}

	// ensure that all Favorites are sanitized
	// to avoid unsafe HTML content
	for i, v := range d.Admins {
		d.Admins[i] = util.Sanitize(v)
	}

	return nil
}

// DashboardFromLibrary converts the library Dashboard type
// to a database Dashboard type.
func FromAPI(d *api.Dashboard) *Dashboard {
	var (
		id  uuid.UUID
		err error
	)

	if d.GetID() == "" {
		id = uuid.New()
	} else {
		id, err = uuid.Parse(d.GetID())
		if err != nil {
			return nil
		}
	}

	user := &Dashboard{
		ID:        id,
		Name:      sql.NullString{String: d.GetName(), Valid: true},
		CreatedAt: sql.NullInt64{Int64: d.GetCreatedAt(), Valid: true},
		CreatedBy: sql.NullString{String: d.GetCreatedBy(), Valid: true},
		UpdatedAt: sql.NullInt64{Int64: d.GetUpdatedAt(), Valid: true},
		UpdatedBy: sql.NullString{String: d.GetUpdatedBy(), Valid: true},
		Admins:    pq.StringArray(d.GetAdmins()),
		Repos:     d.GetRepos(),
	}

	return user.Nullify()
}
