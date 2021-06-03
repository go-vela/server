// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"time"

	"github.com/go-vela/types/constants"
	"github.com/urfave/cli/v2"
)

// Flags represents all supported command line
// interface (CLI) flags for the database.
//
// https://pkg.go.dev/github.com/urfave/cli?tab=doc#Flag
var Flags = []cli.Flag{

	// Logger Flags

	&cli.StringFlag{
		EnvVars: []string{"VELA_DATABASE_LOG_FORMAT", "DATABASE_LOG_FORMAT", "VELA_LOG_FORMAT"},
		Name:    "database.log.format",
		Usage:   "format of logs to output",
		Value:   "json",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_DATABASE_LOG_LEVEL", "DATABASE_LOG_LEVEL", "VELA_LOG_LEVEL"},
		Name:    "database.log.level",
		Usage:   "level of logs to output",
		Value:   "info",
	},

	// Database Flags

	&cli.StringFlag{
		EnvVars: []string{"VELA_DATABASE_DRIVER", "DATABASE_DRIVER"},
		Name:    "database.driver",
		Usage:   "driver to be used for the database system",
		Value:   "sqlite3",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_DATABASE_ADDR", "DATABASE_ADDR"},
		Name:    "database.addr",
		Usage:   "fully qualified url (<scheme>://<host>) for the database",
		Value:   "vela.sqlite",
	},
	&cli.IntFlag{
		EnvVars: []string{"VELA_DATABASE_CONNECTION_OPEN", "DATABASE_CONNECTION_OPEN"},
		Name:    "database.connection.open",
		Usage:   "maximum number of open connections to the database",
		Value:   0,
	},
	&cli.IntFlag{
		EnvVars: []string{"VELA_DATABASE_CONNECTION_IDLE", "DATABASE_CONNECTION_IDLE"},
		Name:    "database.connection.idle",
		Usage:   "maximum number of idle connections to the database",
		Value:   2,
	},
	&cli.DurationFlag{
		EnvVars: []string{"VELA_DATABASE_CONNECTION_LIFE", "DATABASE_CONNECTION_LIFE"},
		Name:    "database.connection.life",
		Usage:   "duration of time a connection may be reused for the database",
		Value:   30 * time.Minute,
	},
	&cli.IntFlag{
		EnvVars: []string{"VELA_DATABASE_COMPRESSION_LEVEL", "DATABASE_COMPRESSION_LEVEL"},
		Name:    "database.compression.level",
		Usage:   "level of compression for logs stored in the database",
		Value:   constants.CompressionThree,
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_DATABASE_ENCRYPTION_KEY", "DATABASE_ENCRYPTION_KEY"},
		Name:    "database.encryption.key",
		Usage:   "AES-256 key for encrypting and decrypting values in the database",
	},
}
