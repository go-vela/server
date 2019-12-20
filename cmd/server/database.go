// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/go-vela/server/database"
	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

// helper function to setup the database from the CLI arguments.
func setupDatabase(c *cli.Context) (database.Service, error) {
	logrus.Debug("Creating database client from CLI configuration")

	switch c.String("database.driver") {
	case constants.DriverPostgres, "postgresql":
		return setupPostgres(c)
	case constants.DriverSqlite, "sqlite":
		return setupSqlite(c)
	default:
		return nil, fmt.Errorf("Unrecognized database driver: %s", c.String("database.driver"))
	}
}

// helper function to setup the Postgres database from the CLI arguments.
func setupPostgres(c *cli.Context) (database.Service, error) {
	logrus.Tracef("Creating %s database client from CLI configuration", constants.DriverPostgres)
	return database.New(c)
}

// helper function to setup the Sqlite database from the CLI arguments.
func setupSqlite(c *cli.Context) (database.Service, error) {
	logrus.Tracef("Creating %s database client from CLI configuration", constants.DriverSqlite)
	return database.New(c)
}
