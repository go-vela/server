// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
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
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableBuildExecutable, err)
	}

	return e, nil
}
