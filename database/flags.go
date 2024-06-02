// SPDX-License-Identifier: Apache-2.0

package database

import (
	"time"

	"github.com/urfave/cli/v2"

	"github.com/go-vela/types/constants"
)

// Flags represents all supported command line interface (CLI) flags for the database.
var Flags = []cli.Flag{
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
	&cli.StringFlag{
		EnvVars:  []string{"VELA_DATABASE_LOG_LEVEL", "DATABASE_LOG_LEVEL"},
		FilePath: "/vela/database/log_level",
		Name:     "database.log.level",
		Usage:    "set log level - options: (trace|info|warn|error)",
		Value:    "warn",
	},
	&cli.BoolFlag{
		EnvVars:  []string{"VELA_DATABASE_LOG_SHOW_SQL", "DATABASE_LOG_SHOW_SQL"},
		FilePath: "/vela/database/log_show_sql",
		Name:     "database.log.show_sql",
		Usage:    "show the SQL query in the logs",
		Value:    false,
	},
	&cli.BoolFlag{
		EnvVars:  []string{"VELA_DATABASE_LOG_SKIP_NOTFOUND", "DATABASE_LOG_SKIP_NOTFOUND"},
		FilePath: "/vela/database/log_skip_notfound",
		Name:     "database.log.skip_notfound",
		Usage:    "skip logging when a resource is not found in the database",
		Value:    true,
	},
	&cli.DurationFlag{
		EnvVars:  []string{"VELA_DATABASE_LOG_SLOW_THRESHOLD", "DATABASE_LOG_SLOW_THRESHOLD"},
		FilePath: "/vela/database/log_slow_threshold",
		Name:     "database.log.slow_threshold",
		Usage:    "queries that take longer than this threshold are considered slow and will be logged",
		Value:    200 * time.Millisecond,
	},
	&cli.BoolFlag{
		EnvVars:  []string{"VELA_DATABASE_SKIP_CREATION", "DATABASE_SKIP_CREATION"},
		FilePath: "/vela/database/skip_creation",
		Name:     "database.skip_creation",
		Usage:    "enables skipping the creation of tables and indexes in the database",
	},
}
