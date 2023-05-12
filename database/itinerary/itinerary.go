// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package itinerary

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type (
	// config represents the settings required to create the engine that implements the BuildItineraryService interface.
	config struct {
		// specifies the level of compression to use for the BuildItinerary engine
		CompressionLevel int
		// specifies to skip creating tables and indexes for the BuildItinerary engine
		SkipCreation bool
		// specifies the driver for proper popping query
		Driver string
	}

	// engine represents the build itinerary functionality that implements the BuildItineraryService interface.
	engine struct {
		// engine configuration settings used in build itinerary functions
		config *config

		// gorm.io/gorm database client used in build itinerary functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in build itinerary functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}
)

// New creates and returns a Vela service for integrating with build itinerariies in the database.
//
//nolint:revive // ignore returning unexported engine
func New(opts ...EngineOpt) (*engine, error) {
	// create new BuildItinerary engine
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

	// check if we should skip creating build itinerary database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of build itineraries table and indexes in the database")

		return e, nil
	}

	// create the build itineraries table
	err := e.CreateBuildItineraryTable(e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableBuildItinerary, err)
	}

	return e, nil
}
