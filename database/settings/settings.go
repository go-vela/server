// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"context"
	"database/sql"
	"fmt"

	api "github.com/go-vela/server/api/types"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

const (
	// todo: constantsTableSettings -> types.constants.TableSettings
	constantsTableSettings = "settings"
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
		ID     sql.NullInt64  `sql:"id"`
		FooNum sql.NullInt64  `sql:"foo_num"`
		FooStr sql.NullString `sql:"foo_str"`
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
		return nil, fmt.Errorf("unable to create %s table: %w", constantsTableSettings, err)
	}

	// todo: need indexes?

	return e, nil
}

// ToAPI converts the Worker type
// to an API Worker type.
func (s *Settings) ToAPI() *api.Settings {
	settings := new(api.Settings)

	// settings.SetID(s.ID.Int64)

	return settings
}
