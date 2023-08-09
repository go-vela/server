// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	"context"
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type (
	// config represents the settings required to create the engine that implements the ScheduleInterface interface.
	config struct {
		// specifies to skip creating tables and indexes for the Schedule engine
		SkipCreation bool
	}

	// engine represents the schedule functionality that implements the ScheduleInterface interface.
	engine struct {
		// engine configuration settings used in schedule functions
		config *config

		ctx context.Context

		// gorm.io/gorm database client used in schedule functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in schedule functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}
)

// New creates and returns a Vela service for integrating with schedules in the database.
//
//nolint:revive // ignore returning unexported engine
func New(opts ...EngineOpt) (*engine, error) {
	// create new Schedule engine
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

	// check if we should skip creating schedule database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of schedules table and indexes in the database")

		return e, nil
	}

	// create the schedules table
	err := e.CreateScheduleTable(e.ctx, e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableSchedule, err)
	}

	// create the indexes for the schedules table
	err = e.CreateScheduleIndexes(e.ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create indexes for %s table: %w", constants.TableSchedule, err)
	}

	return e, nil
}
