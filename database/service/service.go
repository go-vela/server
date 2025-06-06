// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-vela/server/constants"
)

type (
	// config represents the settings required to create the engine that implements the ServiceInterface interface.
	config struct {
		// specifies to skip creating tables and indexes for the Service engine
		SkipCreation bool
	}

	// Engine represents the service functionality that implements the ServiceInterface interface.
	Engine struct {
		// engine configuration settings used in service functions
		config *config

		ctx context.Context

		// gorm.io/gorm database client used in service functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in service functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}
)

// New creates and returns a Vela service for integrating with services in the database.
func New(opts ...EngineOpt) (*Engine, error) {
	// create new Service engine
	e := new(Engine)

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

	// check if we should skip creating service database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of services table in the database")

		return e, nil
	}

	// create the services table
	err := e.CreateServiceTable(e.ctx, e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableService, err)
	}

	return e, nil
}
