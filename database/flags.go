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
		EnvVars:  []string{"VELA_DATABASE_LOG_FORMAT", "DATABASE_LOG_FORMAT", "VELA_LOG_FORMAT"},
		FilePath: "/vela/database/log_format",
		Name:     "database.log.format",
		Usage:    "format of logs to output",
		Value:    "json",
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_DATABASE_LOG_LEVEL", "DATABASE_LOG_LEVEL", "VELA_LOG_LEVEL"},
		FilePath: "/vela/database/log_level",
		Name:     "database.log.level",
		Usage:    "level of logs to output",
		Value:    "info",
	},

	// Database Flags

	&cli.StringFlag{
		EnvVars:  []string{"VELA_DATABASE_DRIVER", "DATABASE_DRIVER"},
		FilePath: "/vela/database/driver",
		Name:     "database.driver",
		Usage:    "driver to be used for the database system",
		Value:    "sqlite3",
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_DATABASE_ADDR", "DATABASE_ADDR"},
		FilePath: "/vela/database/addr",
		Name:     "database.addr",
		Usage:    "fully qualified url (<scheme>://<host>) for the database",
		Value:    "vela.sqlite",
	},
	&cli.IntFlag{
		EnvVars:  []string{"VELA_DATABASE_CONNECTION_OPEN", "DATABASE_CONNECTION_OPEN"},
		FilePath: "/vela/database/connection_open",
		Name:     "database.connection.open",
		Usage:    "maximum number of open connections to the database",
		Value:    0,
	},
	&cli.IntFlag{
		EnvVars:  []string{"VELA_DATABASE_CONNECTION_IDLE", "DATABASE_CONNECTION_IDLE"},
		FilePath: "/vela/database/connection_idle",
		Name:     "database.connection.idle",
		Usage:    "maximum number of idle connections to the database",
		Value:    2,
	},
	&cli.DurationFlag{
		EnvVars:  []string{"VELA_DATABASE_CONNECTION_LIFE", "DATABASE_CONNECTION_LIFE"},
		FilePath: "/vela/database/connection_life",
		Name:     "database.connection.life",
		Usage:    "duration of time a connection may be reused for the database",
		Value:    30 * time.Minute,
	},
	&cli.IntFlag{
		EnvVars:  []string{"VELA_DATABASE_COMPRESSION_LEVEL", "DATABASE_COMPRESSION_LEVEL"},
		FilePath: "/vela/database/compression_level",
		Name:     "database.compression.level",
		Usage:    "level of compression for logs stored in the database",
		Value:    constants.CompressionThree,
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_DATABASE_ENCRYPTION_KEY", "DATABASE_ENCRYPTION_KEY"},
		FilePath: "/vela/database/encryption_key",
		Name:     "database.encryption.key",
		Usage:    "AES-256 key for encrypting and decrypting values in the database",
	},
	&cli.BoolFlag{
		EnvVars:  []string{"VELA_DATABASE_SKIP_CREATION", "DATABASE_SKIP_CREATION"},
		FilePath: "/vela/database/skip_creation",
		Name:     "database.skip_creation",
		Usage:    "enables skipping the creation of objects in the database",
	},
}
