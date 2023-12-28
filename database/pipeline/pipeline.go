// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type (
	// config represents the settings required to create the engine that implements the PipelineInterface interface.
	config struct {
		// specifies the level of compression to use for the Pipeline engine
		CompressionLevel int
		// specifies to skip creating tables and indexes for the Pipeline engine
		SkipCreation bool
	}

	// engine represents the pipeline functionality that implements the PipelineInterface interface.
	engine struct {
		// engine configuration settings used in pipeline functions
		config *config

		ctx context.Context

		// gorm.io/gorm database client used in pipeline functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in pipeline functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}
)

// New creates and returns a Vela service for integrating with pipelines in the database.
//
//nolint:revive // ignore returning unexported engine
func New(opts ...EngineOpt) (*engine, error) {
	// create new Pipeline engine
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

	// check if we should skip creating pipeline database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of pipelines table and indexes in the database")

		return e, nil
	}

	// create the pipelines table
	err := e.CreatePipelineTable(e.ctx, e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TablePipeline, err)
	}

	// create the indexes for the pipelines table
	err = e.CreatePipelineIndexes(e.ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create indexes for %s table: %w", constants.TablePipeline, err)
	}

	return e, nil
}
