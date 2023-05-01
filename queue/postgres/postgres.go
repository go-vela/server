// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"fmt"
	"time"

	"github.com/go-vela/server/queue/postgres/ddl"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type config struct {
	// specifies the address to use for the Postgres client
	Address string
	// specifies a list of channels for managing builds for the Postgres client
	Channels []string
	// specifies the timeout to use for the Postgres client
	Timeout time.Duration
	// specifies the level of compression to use for the Postgres client
	CompressionLevel int
	// specifies the connection duration to use for the Postgres client
	ConnectionLife time.Duration
	// specifies the maximum idle connections for the Postgres client
	ConnectionIdle int
	// specifies the maximum open connections for the Postgres client
	ConnectionOpen int
	// specifies the encryption key to use for the Postgres client
	EncryptionKey string
	// specifies to skip creating tables and indexes for the Postgres client
	SkipCreation bool
}

type client struct {
	config *config
	// https://pkg.go.dev/gorm.io/gorm#DB
	Postgres *gorm.DB
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
	Logger *logrus.Entry
}

// New returns a Queue implementation that
// integrates with a Postgres database instance.
//
//nolint:revive // ignore returning unexported client
func New(opts ...ClientOpt) (*client, error) {
	// create new Postgres client
	c := new(client)

	// create new fields
	c.config = new(config)
	c.Postgres = new(gorm.DB)

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#StandardLogger
	logger := logrus.StandardLogger()

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#NewEntry
	c.Logger = logrus.NewEntry(logger).WithField("queue", c.Driver())

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// create the new Postgres database client
	//
	// https://pkg.go.dev/gorm.io/gorm#Open
	_postgres, err := gorm.Open(postgres.Open(c.config.Address), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// set the Postgres database client in the Postgres client
	c.Postgres = _postgres

	// setup database with proper configuration
	err = setupDatabase(c)
	if err != nil {
		return nil, err
	}

	// // create the services for the database
	// err = createServices(c)
	// if err != nil {
	// 	return nil, err
	// }

	return c, nil
}

// setupDatabase is a helper function to setup
// the database with the proper configuration.
func setupDatabase(c *client) error {
	// capture database/sql database from gorm database
	//
	// https://pkg.go.dev/gorm.io/gorm#DB.DB
	_sql, err := c.Postgres.DB()
	if err != nil {
		return err
	}

	// set the maximum amount of time a connection may be reused
	//
	// https://golang.org/pkg/database/sql/#DB.SetConnMaxLifetime
	_sql.SetConnMaxLifetime(c.config.ConnectionLife)

	// set the maximum number of connections in the idle connection pool
	//
	// https://golang.org/pkg/database/sql/#DB.SetMaxIdleConns
	_sql.SetMaxIdleConns(c.config.ConnectionIdle)

	// set the maximum number of open connections to the database
	//
	// https://golang.org/pkg/database/sql/#DB.SetMaxOpenConns
	_sql.SetMaxOpenConns(c.config.ConnectionOpen)

	// verify connection to the database
	err = c.Ping()
	if err != nil {
		return err
	}

	// check if we should skip creating database objects
	if c.config.SkipCreation {
		c.Logger.Warning("skipping creation of data tables and indexes in the postgres database")

		return nil
	}

	// create the tables in the database
	err = createTables(c)
	if err != nil {
		return err
	}

	// // create the indexes in the database
	// err = createIndexes(c)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// createTables is a helper function to setup
// the database with the necessary tables.
func createTables(c *client) error {
	c.Logger.Trace("creating data tables in the postgres database")

	// create the builds table
	err := c.Postgres.Exec(ddl.CreateBuildsQueueTable).Error
	if err != nil {
		// todo: constants
		return fmt.Errorf("unable to create %s table: %w", BuildsQueueTable, err)
	}

	return nil
}
