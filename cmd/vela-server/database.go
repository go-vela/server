// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/database/postgres"
	"github.com/go-vela/server/database/sqlite"
	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
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
		return nil, fmt.Errorf("invalid database driver: %s", c.String("database.driver"))
	}
}

// helper function to setup the Postgres database from the CLI arguments.
func setupPostgres(c *cli.Context) (database.Service, error) {
	logrus.Tracef("Creating %s database client from CLI configuration", constants.DriverPostgres)

	// create new Postgres database service
	//
	// https://pkg.go.dev/github.com/go-vela/server/database/postgres?tab=doc#New
	return postgres.New(
		postgres.WithAddress(c.String("database.config")),
		postgres.WithCompressionLevel(c.Int("database.compression.level")),
		postgres.WithConnectionLife(c.Duration("database.connection.life")),
		postgres.WithConnectionIdle(c.Int("database.connection.idle")),
		postgres.WithConnectionOpen(c.Int("database.connection.open")),
		postgres.WithEncryptionKey(c.String("database.encryption.key")),
	)
}

// helper function to setup the Sqlite database from the CLI arguments.
func setupSqlite(c *cli.Context) (database.Service, error) {
	logrus.Tracef("Creating %s database client from CLI configuration", constants.DriverSqlite)

	// create new Sqlite database service
	//
	// https://pkg.go.dev/github.com/go-vela/server/database/sqlite?tab=doc#New
	return sqlite.New(
		sqlite.WithAddress(c.String("database.config")),
		sqlite.WithCompressionLevel(c.Int("database.compression.level")),
		sqlite.WithConnectionLife(c.Duration("database.connection.life")),
		sqlite.WithConnectionIdle(c.Int("database.connection.idle")),
		sqlite.WithConnectionOpen(c.Int("database.connection.open")),
		sqlite.WithEncryptionKey(c.String("database.encryption.key")),
	)
}
