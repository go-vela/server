// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	TableSettings = "settings"
)

type (
	// config represents the settings required to create the engine that implements the SettingsInterface interface.
	config struct {
		// specifies to skip creating tables and indexes for the Settings engine
		SkipCreation bool
	}

	// engine represents the settings functionality that implements the SettingsInterface interface.
	engine struct {
		// engine configuration settings used in settings functions
		config *config

		ctx context.Context

		// gorm.io/gorm database client used in settings functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in settings functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}
)

// New creates and returns a Vela service for integrating with settings in the database.
//
//nolint:revive // ignore returning unexported engine
func New(opts ...EngineOpt) (*engine, error) {
	// create new Settings engine
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
		e.logger.Warning("skipping creation of settings table and indexes")

		return e, nil
	}

	// create the settings table
	err := e.CreateSettingsTable(e.ctx, e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", TableSettings, err)
	}

	return e, nil
}
