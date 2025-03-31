// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-vela/server/constants"
)

type (
	// config represents the settings required to create the engine that implements the LogInterface interface.
	config struct {
		// specifies the level of compression to use for the Log engine
		CompressionLevel int
		// specifies to skip creating tables and indexes for the Log engine
		SkipCreation bool
	}

	// engine represents the log functionality that implements the LogInterface interface.
	engine struct {
		// engine configuration settings used in log functions
		config *config

		ctx context.Context

		// gorm.io/gorm database client used in log functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in log functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}
)

// New creates and returns a Vela service for integrating with logs in the database.
//
//nolint:revive // ignore returning unexported client
func New(opts ...EngineOpt) (*engine, error) {
	// create new Log engine
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

	// check if we should skip creating log database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of logs table and indexes")

		return e, nil
	}

	// create the logs table
	err := e.CreateLogTable(e.ctx, e.client.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableLog, err)
	}

	// create the indexes for the logs table
	err = e.CreateLogIndexes(e.ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create indexes for %s table: %w", constants.TableLog, err)
	}

	return e, nil
}
