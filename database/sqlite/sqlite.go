// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"fmt"
	"time"

	"github.com/go-vela/server/database/sqlite/ddl"
	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type (
	config struct {
		// specifies the address to use for the Sqlite client
		Address string
		// specifies the level of compression to use for the Sqlite client
		CompressionLevel int
		// specifies the connection duration to use for the Sqlite client
		ConnectionLife time.Duration
		// specifies the maximum idle connections for the Sqlite client
		ConnectionIdle int
		// specifies the maximum open connections for the Sqlite client
		ConnectionOpen int
		// specifies the encryption key to use for the Sqlite client
		EncryptionKey string
	}

	client struct {
		config *config
		Sqlite *gorm.DB
	}
)

// New returns a Database implementation that integrates with a Sqlite instance.
//
// nolint: golint // ignore returning unexported client
func New(opts ...ClientOpt) (*client, error) {
	// create new Sqlite client
	c := new(client)

	// create new fields
	c.config = new(config)
	c.Sqlite = new(gorm.DB)

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// create the new Sqlite database client
	//
	// https://pkg.go.dev/gorm.io/gorm#Open
	_sqlite, err := gorm.Open(sqlite.Open(c.config.Address), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// set the Sqlite database client in the Sqlite client
	c.Sqlite = _sqlite

	// setup database with proper configuration
	err = setupDatabase(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// NewTest returns a Database implementation that integrates with a fake Sqlite instance.
//
// This function is intended for running tests only.
//
// nolint: golint // ignore returning unexported client
func NewTest() (*client, error) {
	// create new Sqlite client
	c := new(client)

	// create new fields
	c.config = &config{
		CompressionLevel: 3,
		ConnectionLife:   30 * time.Minute,
		ConnectionIdle:   2,
		ConnectionOpen:   0,
		EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
	}
	c.Sqlite = new(gorm.DB)

	// create the new Sqlite database client
	//
	// https://pkg.go.dev/gorm.io/gorm#Open
	_sqlite, err := gorm.Open(
		sqlite.Open("file::memory:?cache=shared"),
		&gorm.Config{SkipDefaultTransaction: true},
	)
	if err != nil {
		return nil, err
	}

	c.Sqlite = _sqlite

	return c, nil
}

// setupDatabase is a helper function to setup
// the database with the proper configuration.
func setupDatabase(c *client) error {
	// capture database/sql database from gorm database
	//
	// https://pkg.go.dev/gorm.io/gorm#DB.DB
	_sql, err := c.Sqlite.DB()
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

	// create the tables in the database
	err = createTables(c)
	if err != nil {
		return err
	}

	return nil
}

// createTables is a helper function to setup
// the database with the necessary tables.
func createTables(c *client) error {
	logrus.Trace("creating data tables in the sqlite database")

	// create the builds table
	err := c.Sqlite.Exec(ddl.CreateBuildTable).Error
	if err != nil {
		return fmt.Errorf("unable to create %s table: %v", constants.TableBuild, err)
	}

	// create the hooks table
	err = c.Sqlite.Exec(ddl.CreateHookTable).Error
	if err != nil {
		return fmt.Errorf("unable to create %s table: %v", constants.TableHook, err)
	}

	// create the logs table
	err = c.Sqlite.Exec(ddl.CreateLogTable).Error
	if err != nil {
		return fmt.Errorf("unable to create %s table: %v", constants.TableLog, err)
	}

	// create the repos table
	err = c.Sqlite.Exec(ddl.CreateRepoTable).Error
	if err != nil {
		return fmt.Errorf("unable to create %s table: %v", constants.TableRepo, err)
	}

	// create the secrets table
	err = c.Sqlite.Exec(ddl.CreateSecretTable).Error
	if err != nil {
		return fmt.Errorf("unable to create %s table: %v", constants.TableSecret, err)
	}

	// create the services table
	err = c.Sqlite.Exec(ddl.CreateServiceTable).Error
	if err != nil {
		return fmt.Errorf("unable to create %s table: %v", constants.TableService, err)
	}

	// create the steps table
	err = c.Sqlite.Exec(ddl.CreateStepTable).Error
	if err != nil {
		return fmt.Errorf("unable to create %s table: %v", constants.TableStep, err)
	}

	// create the users table
	err = c.Sqlite.Exec(ddl.CreateUserTable).Error
	if err != nil {
		return fmt.Errorf("unable to create %s table: %v", constants.TableUser, err)
	}

	// create the workers table
	err = c.Sqlite.Exec(ddl.CreateWorkerTable).Error
	if err != nil {
		return fmt.Errorf("unable to create %s table: %v", constants.TableWorker, err)
	}

	return nil
}
