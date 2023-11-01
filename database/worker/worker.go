// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type (
	// config represents the settings required to create the engine that implements the WorkerInterface interface.
	config struct {
		// specifies to skip creating tables and indexes for the Worker engine
		SkipCreation bool
	}

	// engine represents the worker functionality that implements the WorkerInterface interface.
	engine struct {
		// engine configuration settings used in worker functions
		config *config

		ctx context.Context

		// gorm.io/gorm database client used in worker functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in worker functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}
)

// New creates and returns a Vela service for integrating with workers in the database.
//
//nolint:revive // ignore returning unexported engine
func New(opts ...EngineOpt) (*engine, error) {
	// create new Worker engine
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

	// check if we should skip creating worker database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of workers table and indexes in the database")

		return e, nil
	}

	// create the workers table
	err := e.CreateWorkerTable(e.ctx, e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableWorker, err)
	}

	// create the indexes for the workers table
	err = e.CreateWorkerIndexes(e.ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create indexes for %s table: %w", constants.TableWorker, err)
	}

	return e, nil
}

// convertToBuilds is a helper function that generates build objects with ID fields given a list of IDs.
func convertToBuilds(ids []string) []*library.Build {
	// create stripped build objects holding the IDs
	var rBs []*library.Build

	for _, b := range ids {
		id, err := strconv.ParseInt(b, 10, 64)
		if err != nil {
			return nil
		}

		build := new(library.Build)
		build.SetID(id)

		rBs = append(rBs, build)
	}

	return rBs
}
