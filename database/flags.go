// SPDX-License-Identifier: Apache-2.0

package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/go-vela/server/constants"
)

// Flags represents all supported command line interface (CLI) flags for the database.
var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:  "database.driver",
		Usage: "driver to be used for the database system",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_DATABASE_DRIVER"),
			cli.EnvVar("DATABASE_DRIVER"),
			cli.File("/vela/database/driver"),
		),
		Value: "sqlite3",
		Action: func(_ context.Context, _ *cli.Command, v string) error {
			if len(v) == 0 {
				return fmt.Errorf("no database driver provided")
			}

			return nil
		},
	},
	&cli.StringFlag{
		Name:  "database.addr",
		Usage: "fully qualified url (<scheme>://<host>) for the database",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_DATABASE_ADDR"),
			cli.EnvVar("DATABASE_ADDR"),
			cli.File("/vela/database/addr"),
		),
		Required: true,
		Action: func(_ context.Context, _ *cli.Command, v string) error {
			if strings.HasSuffix(v, "/") {
				return fmt.Errorf("invalid database address provided: address must not have trailing slash")
			}

			return nil
		},
	},
	&cli.IntFlag{
		Name:  "database.connection.open",
		Usage: "maximum number of open connections to the database",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_DATABASE_CONNECTION_OPEN"),
			cli.EnvVar("DATABASE_CONNECTION_OPEN"),
			cli.File("/vela/database/connection_open"),
		),
		Value: 0,
	},
	&cli.IntFlag{
		Name:  "database.connection.idle",
		Usage: "maximum number of idle connections to the database",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_DATABASE_CONNECTION_IDLE"),
			cli.EnvVar("DATABASE_CONNECTION_IDLE"),
			cli.File("/vela/database/connection_idle"),
		),
		Value: 2,
	},
	&cli.DurationFlag{
		Name:  "database.connection.life",
		Usage: "duration of time a connection may be reused for the database",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_DATABASE_CONNECTION_LIFE"),
			cli.EnvVar("DATABASE_CONNECTION_LIFE"),
			cli.File("/vela/database/connection_life"),
		),
		Value: 30 * time.Minute,
	},
	&cli.Int64Flag{
		Name:  "database.compression.level",
		Usage: "level of compression for logs stored in the database",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_DATABASE_COMPRESSION_LEVEL"),
			cli.EnvVar("DATABASE_COMPRESSION_LEVEL"),
			cli.File("/vela/database/compression_level"),
		),
		Value: constants.CompressionThree,
		Action: func(_ context.Context, _ *cli.Command, v int64) error {
			switch v {
			case constants.CompressionNegOne:
				fallthrough
			case constants.CompressionZero:
				fallthrough
			case constants.CompressionOne:
				fallthrough
			case constants.CompressionTwo:
				fallthrough
			case constants.CompressionThree:
				fallthrough
			case constants.CompressionFour:
				fallthrough
			case constants.CompressionFive:
				fallthrough
			case constants.CompressionSix:
				fallthrough
			case constants.CompressionSeven:
				fallthrough
			case constants.CompressionEight:
				fallthrough
			case constants.CompressionNine:
				break
			default:
				return fmt.Errorf("invalid database compression level provided: level (%d) must be between %d and %d",
					v, constants.CompressionNegOne, constants.CompressionNine,
				)
			}

			return nil
		},
	},
	&cli.StringFlag{
		Name:  "database.encryption.key",
		Usage: "AES-256 key for encrypting and decrypting values in the database",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_DATABASE_ENCRYPTION_KEY"),
			cli.EnvVar("DATABASE_ENCRYPTION_KEY"),
			cli.File("/vela/database/encryption_key"),
		),
		Required: true,
		Action: func(_ context.Context, _ *cli.Command, v string) error {
			if len(v) != 32 {
				return fmt.Errorf("invalid database encryption key provided: key length (%d) must be 32 characters", len(v))
			}

			return nil
		},
	},
	&cli.StringFlag{
		Name:  "database.log.level",
		Usage: "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_DATABASE_LOG_LEVEL"),
			cli.EnvVar("DATABASE_LOG_LEVEL"),
			cli.File("/vela/database/log_level"),
		),
		Value: "warn",
	},
	&cli.BoolFlag{
		Name:  "database.log.show_sql",
		Usage: "show the SQL query in the logs",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_DATABASE_LOG_SHOW_SQL"),
			cli.EnvVar("DATABASE_LOG_SHOW_SQL"),
			cli.File("/vela/database/log_show_sql"),
		),
	},
	&cli.BoolFlag{
		Name:  "database.log.skip_notfound",
		Usage: "skip logging when a resource is not found in the database",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_DATABASE_LOG_SKIP_NOTFOUND"),
			cli.EnvVar("DATABASE_LOG_SKIP_NOTFOUND"),
			cli.File("/vela/database/log_skip_notfound"),
		),
		Value: true,
	},
	&cli.DurationFlag{
		Name:  "database.log.slow_threshold",
		Usage: "queries that take longer than this threshold are considered slow and will be logged",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_DATABASE_LOG_SLOW_THRESHOLD"),
			cli.EnvVar("DATABASE_LOG_SLOW_THRESHOLD"),
			cli.File("/vela/database/log_slow_threshold"),
		),
		Value: 200 * time.Millisecond,
	},
	&cli.BoolFlag{
		Name:  "database.skip_creation",
		Usage: "enables skipping the creation of tables and indexes in the database",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_DATABASE_SKIP_CREATION"),
			cli.EnvVar("DATABASE_SKIP_CREATION"),
			cli.File("/vela/database/skip_creation"),
		),
	},
	&cli.BoolFlag{
		Name:  "database.log.partitioned",
		Usage: "enables partition-aware log cleanup for partitioned log tables",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_DATABASE_LOG_PARTITIONED"),
			cli.EnvVar("DATABASE_LOG_PARTITIONED"),
			cli.File("/vela/database/log_partitioned"),
		),
		Value: false,
	},
	&cli.StringFlag{
		Name:  "database.log.partition_pattern",
		Usage: "naming pattern for log table partitions (e.g., logs_%, logs_y%, logs_monthly_)",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_DATABASE_LOG_PARTITION_PATTERN"),
			cli.EnvVar("DATABASE_LOG_PARTITION_PATTERN"),
			cli.File("/vela/database/log_partition_pattern"),
		),
		Value: "logs_%",
	},
	&cli.StringFlag{
		Name:  "database.log.partition_schema",
		Usage: "database schema containing log table partitions",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_DATABASE_LOG_PARTITION_SCHEMA"),
			cli.EnvVar("DATABASE_LOG_PARTITION_SCHEMA"),
			cli.File("/vela/database/log_partition_schema"),
		),
		Value: "public",
	},
}
