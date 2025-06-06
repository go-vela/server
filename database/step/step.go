// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-vela/server/constants"
)

type (
	// config represents the settings required to create the engine that implements the StepInterface interface.
	config struct {
		// specifies to skip creating tables and indexes for the Step engine
		SkipCreation bool
	}

	// Engine represents the step functionality that implements the StepInterface interface.
	Engine struct {
		// engine configuration settings used in step functions
		config *config

		ctx context.Context

		// gorm.io/gorm database client used in step functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in step functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}
)

// New creates and returns a Vela service for integrating with steps in the database.
func New(opts ...EngineOpt) (*Engine, error) {
	// create new Step engine
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

	// check if we should skip creating step database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of steps table in the database")

		return e, nil
	}

	// create the steps table
	err := e.CreateStepTable(e.ctx, e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableStep, err)
	}

	return e, nil
}
