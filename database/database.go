// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"fmt"
	"time"

	"github.com/go-vela/server/database/build"
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
	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type (
	// Config represents the settings required to create the engine that implements the Interface.
	Config struct {
		// specifies the address to use for the database client
		Address string
		// specifies the level of compression to use for the database client
		CompressionLevel int
		// specifies the connection duration to use for the database client
		ConnectionLife time.Duration
		// specifies the maximum idle connections for the database client
		ConnectionIdle int
		// specifies the maximum open connections for the database client
		ConnectionOpen int
		// specifies the driver to use for the database client
		Driver string
		// specifies the encryption key to use for the database client
		EncryptionKey string
		// specifies to skip creating tables and indexes for the database client
		SkipCreation bool
	}

	// engine represents the functionality that implements the Interface.
	engine struct {
		Config   *Config
		Database *gorm.DB
		Logger   *logrus.Entry

		build.BuildInterface
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
func New(c *Config) (Interface, error) {
	// validate the configuration being provided
	err := c.Validate()
	if err != nil {
		return nil, err
	}

	// create new database engine
	e := &engine{
		Config:   c,
		Database: new(gorm.DB),
		Logger:   logrus.NewEntry(logrus.StandardLogger()).WithField("database", c.Driver),
	}

	e.Logger.Trace("creating database engine from configuration")
	// process the database driver being provided
	switch c.Driver {
	case constants.DriverPostgres:
		// create the new Postgres database client
		e.Database, err = gorm.Open(postgres.Open(e.Config.Address), &gorm.Config{})
		if err != nil {
			return nil, err
		}
	case constants.DriverSqlite:
		// create the new Sqlite database client
		e.Database, err = gorm.Open(sqlite.Open(e.Config.Address), &gorm.Config{})
		if err != nil {
			return nil, err
		}
	default:
		// handle an invalid database driver being provided
		return nil, fmt.Errorf("invalid database driver provided: %s", c.Driver)
	}

	// capture database/sql database from gorm.io/gorm database
	db, err := e.Database.DB()
	if err != nil {
		return nil, err
	}

	// set the maximum amount of time a connection may be reused
	db.SetConnMaxLifetime(e.Config.ConnectionLife)
	// set the maximum number of connections in the idle connection pool
	db.SetMaxIdleConns(e.Config.ConnectionIdle)
	// set the maximum number of open connections to the database
	db.SetMaxOpenConns(e.Config.ConnectionOpen)

	// verify connection to the database
	err = e.Ping()
	if err != nil {
		return nil, err
	}

	// create database agnostic engines for resources
	err = e.NewResources()
	if err != nil {
		return nil, err
	}

	return e, nil
}

// NewTest creates and returns an engine that integrates with an in-memory database provider.
//
// This function is ONLY intended to be used for testing purposes.
func NewTest() (Interface, error) {
	return New(&Config{
		Address:          "file::memory:?cache=shared",
		CompressionLevel: 3,
		ConnectionLife:   30 * time.Minute,
		ConnectionIdle:   2,
		ConnectionOpen:   0,
		Driver:           "sqlite3",
		EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
		SkipCreation:     false,
	})
}
