// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-vela/server/database/build"
	"github.com/go-vela/server/database/executable"
	"github.com/go-vela/server/database/hook"
	"github.com/go-vela/server/database/log"
	"github.com/go-vela/server/database/pipeline"
	"github.com/go-vela/server/database/repo"
	"github.com/go-vela/server/database/schedule"
	"github.com/go-vela/server/database/secret"
	"github.com/go-vela/server/database/service"
	"github.com/go-vela/server/database/step"
	"github.com/go-vela/server/database/user"
	"github.com/go-vela/server/database/worker"
	"github.com/go-vela/server/tracing"
	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type (
	// config represents the settings required to create the engine that implements the Interface.
	config struct {
		// specifies the address to use for the database engine
		Address string
		// specifies the level of compression to use for the database engine
		CompressionLevel int
		// specifies the connection duration to use for the database engine
		ConnectionLife time.Duration
		// specifies the maximum idle connections for the database engine
		ConnectionIdle int
		// specifies the maximum open connections for the database engine
		ConnectionOpen int
		// specifies the driver to use for the database engine
		Driver string
		// specifies the encryption key to use for the database engine
		EncryptionKey string
		// specifies to skip creating tables and indexes for the database engine
		SkipCreation bool
	}

	// engine represents the functionality that implements the Interface.
	engine struct {
		// gorm.io/gorm database client used in database functions
		client *gorm.DB
		// engine configuration settings used in database functions
		config *config
		// engine context used in database functions
		ctx context.Context
		// sirupsen/logrus logger used in database functions
		logger *logrus.Entry
		// configurations related to telemetry/tracing
		tracing *tracing.Config

		build.BuildInterface
		executable.BuildExecutableInterface
		hook.HookInterface
		log.LogInterface
		pipeline.PipelineInterface
		repo.RepoInterface
		schedule.ScheduleInterface
		secret.SecretInterface
		service.ServiceInterface
		step.StepInterface
		user.UserInterface
		worker.WorkerInterface
	}
)

// New creates and returns an engine capable of integrating with the configured database provider.
//
// Currently, the following database providers are supported:
//
// * postgres
// * sqlite3
func New(opts ...EngineOpt) (Interface, error) {
	// create new database engine
	e := new(engine)

	// create new fields
	e.client = new(gorm.DB)
	e.config = new(config)
	e.logger = new(logrus.Entry)
	e.ctx = context.TODO()

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(e)
		if err != nil {
			return nil, err
		}
	}

	// validate the configuration being provided
	err := e.config.Validate()
	if err != nil {
		return nil, err
	}

	// update the logger with additional metadata
	e.logger = logrus.NewEntry(logrus.StandardLogger()).WithField("database", e.Driver())

	e.logger.Trace("creating database engine from configuration")
	// process the database driver being provided
	switch e.config.Driver {
	case constants.DriverPostgres:
		// create the new Postgres database client
		e.client, err = gorm.Open(postgres.Open(e.config.Address), &gorm.Config{})
		if err != nil {
			return nil, err
		}
	case constants.DriverSqlite:
		// create the new Sqlite database client
		e.client, err = gorm.Open(sqlite.Open(e.config.Address), &gorm.Config{})
		if err != nil {
			return nil, err
		}
	default:
		// handle an invalid database driver being provided
		return nil, fmt.Errorf("invalid database driver provided: %s", e.Driver())
	}

	// capture database/sql database from gorm.io/gorm database
	db, err := e.client.DB()
	if err != nil {
		return nil, err
	}

	// initialize otel tracing if enabled
	if e.tracing.EnableTracing {
		otelPlugin := otelgorm.NewPlugin(
			otelgorm.WithTracerProvider(e.tracing.TracerProvider),
			otelgorm.WithoutQueryVariables(),
		)

		err := e.client.Use(otelPlugin)
		if err != nil {
			return nil, err
		}
	}

	// set the maximum amount of time a connection may be reused
	db.SetConnMaxLifetime(e.config.ConnectionLife)
	// set the maximum number of connections in the idle connection pool
	db.SetMaxIdleConns(e.config.ConnectionIdle)
	// set the maximum number of open connections to the database
	db.SetMaxOpenConns(e.config.ConnectionOpen)

	// verify connection to the database
	err = e.Ping()
	if err != nil {
		return nil, err
	}

	// create database agnostic engines for resources
	err = e.NewResources(e.ctx)
	if err != nil {
		return nil, err
	}

	return e, nil
}

// NewTest creates and returns an engine that integrates with an in-memory database provider.
//
// This function is ONLY intended to be used for testing purposes.
func NewTest() (Interface, error) {
	return New(
		WithAddress("file::memory:?cache=shared"),
		WithCompressionLevel(3),
		WithConnectionLife(30*time.Minute),
		WithConnectionIdle(2),
		WithConnectionOpen(0),
		WithDriver("sqlite3"),
		WithEncryptionKey("A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW"),
		WithSkipCreation(false),
		WithTracingConfig(&tracing.Config{EnableTracing: false}),
	)
}
