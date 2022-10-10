// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package log

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type (
	// config represents the settings required to create the engine that implements the LogService interface.
	config struct {
		// specifies the level of compression to use for the Log engine
		CompressionLevel int
		// specifies to skip creating tables and indexes for the Log engine
		SkipCreation bool
	}

	// engine represents the log functionality that implements the LogService interface.
	engine struct {
		// engine configuration settings used in log functions
		config *config

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
//nolint:revive // ignore returning unexported engine
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
		e.logger.Warning("skipping creation of logs table and indexes in the database")

		return e, nil
	}

	// create the logs table
	err := e.CreateLogTable(e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableLog, err)
	}

	// create the indexes for the logs table
	err = e.CreateLogIndexes()
	if err != nil {
		return nil, fmt.Errorf("unable to create indexes for %s table: %w", constants.TableLog, err)
	}

	return e, nil
}
