// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/go-vela/server/database"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

// helper function to setup the database from the CLI arguments.
func setupDatabase(c *cli.Context) (database.Interface, error) {
	logrus.Debug("Creating database client from CLI configuration")

	// database configuration
	_setup := &database.Setup{
		Driver:           c.String("database.driver"),
		Address:          c.String("database.addr"),
		CompressionLevel: c.Int("database.compression.level"),
		ConnectionLife:   c.Duration("database.connection.life"),
		ConnectionIdle:   c.Int("database.connection.idle"),
		ConnectionOpen:   c.Int("database.connection.open"),
		EncryptionKey:    c.String("database.encryption.key"),
		SkipCreation:     c.Bool("database.skip_creation"),
	}

	// setup the database
	//
	// https://pkg.go.dev/github.com/go-vela/server/database?tab=doc#New
	return database.New(_setup)
}
