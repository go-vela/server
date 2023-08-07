// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

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

		// gorm.io/gorm database client used in build functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in build functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}
)

// New creates and returns a Vela service for integrating with builds in the database.
func New(opts ...EngineOpt) (BuildInterface, error) {
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

	// check if we should skip creating build database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of builds table and indexes in the database")

		return e, nil
	}

	// create the builds table
	err := e.CreateBuildTable(e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableBuild, err)
	}

	// create the indexes for the builds table
	err = e.CreateBuildIndexes()
	if err != nil {
		return nil, fmt.Errorf("unable to create indexes for %s table: %w", constants.TableBuild, err)
	}

	return e, nil
}
