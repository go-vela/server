// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/go-vela/server/database/ddl"
	"github.com/go-vela/server/database/dml"

	"github.com/go-vela/types/constants"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type client struct {
	Database *gorm.DB

	DDL *ddl.Map
	DML *dml.Map
}

// New returns a Database implementation that
// integrates with a supported database instance.
func New(c *cli.Context) (*client, error) {
	driver := c.String("database.driver")
	config := c.String("database.config")

	// create the database client
	db, err := gorm.Open(driver, config)
	if err != nil {
		return nil, err
	}

	// create the DDL map
	ddlMap, err := ddl.NewMap(driver)
	if err != nil {
		return nil, err
	}

	// create the DML map
	dmlMap, err := dml.NewMap(driver)
	if err != nil {
		return nil, err
	}

	// setup database with proper configuration
	err = setupDatabase(db.DB(), ddlMap)
	if err != nil {
		return nil, err
	}

	// apply extra database configuration
	db.DB().SetConnMaxLifetime(c.Duration("database.connection.life"))
	db.DB().SetMaxIdleConns(c.Int("database.connection.idle"))
	db.DB().SetMaxOpenConns(c.Int("database.connection.open"))

	// create the client object
	client := &client{
		Database: db,
		DDL:      ddlMap,
		DML:      dmlMap,
	}

	return client, nil
}

// NewTest returns a Database implementation that
// integrates with an in-memory Sqlite database instance.
//
// It's possible to overide this with env variables,
// which gets used as a part of integration testing
// with the different supported backends.
//
// This function is intended for running tests only.
func NewTest() (*client, error) {
	name := os.Getenv("VELA_DATABASE_DRIVER")
	if len(name) == 0 {
		name = constants.DriverSqlite
	}

	config := os.Getenv("VELA_DATABASE_CONFIG")
	if len(config) == 0 {
		config = ":memory:"
	}

	// create the database client
	db, err := gorm.Open(name, config)
	if err != nil {
		return nil, err
	}

	// create the DDL map
	ddlMap, err := ddl.NewMap(name)
	if err != nil {
		return nil, err
	}

	// create the DML map
	dmlMap, err := dml.NewMap(name)
	if err != nil {
		return nil, err
	}

	// since sqlite database is in memory, we
	// need to create the tables everytime
	if name == constants.DriverSqlite {
		err = createTables(db.DB(), ddlMap)
		if err != nil {
			return nil, err
		}
	}

	// create the client object
	client := &client{
		Database: db,
		DDL:      ddlMap,
		DML:      dmlMap,
	}

	return client, nil
}

// setupDatabase is a helper function to setup
// the database with the proper configuration.
func setupDatabase(db *sql.DB, ddlMap *ddl.Map) error {
	// ping the database
	err := pingDatabase(db)
	if err != nil {
		return err
	}

	// create the tables in the database
	err = createTables(db, ddlMap)
	if err != nil {
		return err
	}

	// create the indexes in the database
	err = createIndexes(db, ddlMap)
	if err != nil {
		return err
	}

	return nil
}

// pingDatabase is a helper function to send a
// "ping" request with backoff to the database.
//
// This will ensure we have properly established a
// connection to the database instance before we
// try to set it up.
func pingDatabase(db *sql.DB) error {
	// attempt 10 times
	for i := 0; i < 10; i++ {
		// send ping request to database
		err := db.Ping()
		if err != nil {
			logrus.Debugf("unable to ping database. Retrying in %v", (time.Duration(i) * time.Second))
			time.Sleep(1 * time.Second)

			continue
		}

		return nil
	}

	return fmt.Errorf("unable to establish database connection")
}

// createTables is a helper function to setup
// the database with the necessary tables.
func createTables(db *sql.DB, ddlMap *ddl.Map) error {
	// create the build table
	_, err := db.Exec(ddlMap.BuildService.Create)
	if err != nil {
		return fmt.Errorf("unable to create %s table: %v", constants.TableBuild, err)
	}

	// create the hook table
	_, err = db.Exec(ddlMap.HookService.Create)
	if err != nil {
		return fmt.Errorf("unable to create %s table: %v", constants.TableHook, err)
	}

	// create the log table
	_, err = db.Exec(ddlMap.LogService.Create)
	if err != nil {
		return fmt.Errorf("unable to create %s table: %v", constants.TableLog, err)
	}

	// create the repo table
	_, err = db.Exec(ddlMap.RepoService.Create)
	if err != nil {
		return fmt.Errorf("unable to create %s table: %v", constants.TableRepo, err)
	}

	// create the secret table
	_, err = db.Exec(ddlMap.SecretService.Create)
	if err != nil {
		return fmt.Errorf("unable to create %s table: %v", constants.TableSecret, err)
	}

	// create the step table
	_, err = db.Exec(ddlMap.StepService.Create)
	if err != nil {
		return fmt.Errorf("unable to create %s table: %v", constants.TableStep, err)
	}

	// create the service table
	_, err = db.Exec(ddlMap.ServiceService.Create)
	if err != nil {
		return fmt.Errorf("unable to create %s table: %v", constants.TableService, err)
	}

	// create the user table
	_, err = db.Exec(ddlMap.UserService.Create)
	if err != nil {
		return fmt.Errorf("unable to create %s table: %v", constants.TableUser, err)
	}

	// create the worker table
	_, err = db.Exec(ddlMap.WorkerService.Create)
	if err != nil {
		return fmt.Errorf("unable to create %s table: %v", constants.TableWorker, err)
	}

	return nil
}

// createIndexes is a helper function to setup
// the database with the necessary indexes.
func createIndexes(db *sql.DB, ddlMap *ddl.Map) error {
	// create the build table indexes
	for _, index := range ddlMap.BuildService.Indexes {
		_, err := db.Exec(index)
		if err != nil {
			return fmt.Errorf("unable to create %s table indexes: %v", constants.TableBuild, err)
		}
	}

	// create the hook table indexes
	for _, index := range ddlMap.HookService.Indexes {
		_, err := db.Exec(index)
		if err != nil {
			return fmt.Errorf("unable to create %s table indexes: %v", constants.TableHook, err)
		}
	}

	// create the log table indexes
	for _, index := range ddlMap.LogService.Indexes {
		_, err := db.Exec(index)
		if err != nil {
			return fmt.Errorf("unable to create %s table indexes: %v", constants.TableLog, err)
		}
	}

	// create the repo table indexes
	for _, index := range ddlMap.RepoService.Indexes {
		_, err := db.Exec(index)
		if err != nil {
			return fmt.Errorf("unable to create %s table indexes: %v", constants.TableRepo, err)
		}
	}

	// create the secret table indexes
	for _, index := range ddlMap.SecretService.Indexes {
		_, err := db.Exec(index)
		if err != nil {
			return fmt.Errorf("unable to create %s table indexes: %v", constants.TableSecret, err)
		}
	}

	// create the step table indexes
	for _, index := range ddlMap.StepService.Indexes {
		_, err := db.Exec(index)
		if err != nil {
			return fmt.Errorf("unable to create %s table indexes: %v", constants.TableStep, err)
		}
	}

	// create the service table indexes
	for _, index := range ddlMap.ServiceService.Indexes {
		_, err := db.Exec(index)
		if err != nil {
			return fmt.Errorf("unable to create %s table indexes: %v", constants.TableService, err)
		}
	}

	// create the user table indexes
	for _, index := range ddlMap.UserService.Indexes {
		_, err := db.Exec(index)
		if err != nil {
			return fmt.Errorf("unable to create %s table indexes: %v", constants.TableUser, err)
		}
	}

	// create the worker table indexes
	for _, index := range ddlMap.WorkerService.Indexes {
		_, err := db.Exec(index)
		if err != nil {
			return fmt.Errorf("unable to create %s table indexes: %v", constants.TableWorker, err)
		}
	}

	return nil
}
