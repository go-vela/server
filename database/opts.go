// SPDX-License-Identifier: Apache-2.0

package database

import (
	"context"
	"time"

	"github.com/go-vela/server/tracing"
)

// EngineOpt represents a configuration option to initialize the database engine.
type EngineOpt func(*engine) error

// WithAddress sets the address in the database engine.
func WithAddress(address string) EngineOpt {
	return func(e *engine) error {
		// set the fully qualified connection string in the database engine
		e.config.Address = address

		return nil
	}
}

// WithCompressionLevel sets the compression level in the database engine.
func WithCompressionLevel(level int) EngineOpt {
	return func(e *engine) error {
		// set the level of compression for resources in the database engine
		e.config.CompressionLevel = level

		return nil
	}
}

// WithConnectionLife sets the life of connections in the database engine.
func WithConnectionLife(connectionLife time.Duration) EngineOpt {
	return func(e *engine) error {
		// set the maximum duration of time for connection in the database engine
		e.config.ConnectionLife = connectionLife

		return nil
	}
}

// WithConnectionIdle sets the idle connections in the database engine.
func WithConnectionIdle(connectionIdle int) EngineOpt {
	return func(e *engine) error {
		// set the maximum allowed idle connections in the database engine
		e.config.ConnectionIdle = connectionIdle

		return nil
	}
}

// WithConnectionOpen sets the open connections in the database engine.
func WithConnectionOpen(connectionOpen int) EngineOpt {
	return func(e *engine) error {
		// set the maximum allowed open connections in the database engine
		e.config.ConnectionOpen = connectionOpen

		return nil
	}
}

// WithDriver sets the driver in the database engine.
func WithDriver(driver string) EngineOpt {
	return func(e *engine) error {
		// set the database type to interact with in the database engine
		e.config.Driver = driver

		return nil
	}
}

// WithEncryptionKey sets the encryption key in the database engine.
func WithEncryptionKey(encryptionKey string) EngineOpt {
	return func(e *engine) error {
		// set the key for encrypting resources in the database engine
		e.config.EncryptionKey = encryptionKey

		return nil
	}
}

// WithLogLevel sets the log level in the database engine.
func WithLogLevel(logLevel string) EngineOpt {
	return func(e *engine) error {
		// set the log level for the database engine
		e.config.LogLevel = logLevel

		return nil
	}
}

// WithLogSkipNotFound sets the log skip not found option in the database engine.
func WithLogSkipNotFound(logSkipNotFound bool) EngineOpt {
	return func(e *engine) error {
		// set the log skip not found option for the database engine
		e.config.LogSkipNotFound = logSkipNotFound

		return nil
	}
}

// WithLogSlowThreshold sets the log slow query threshold in the database engine.
func WithLogSlowThreshold(logSlowThreshold time.Duration) EngineOpt {
	return func(e *engine) error {
		// set the slow query threshold for the database engine
		e.config.LogSlowThreshold = logSlowThreshold

		return nil
	}
}

// WithLogShowSQL sets the log show SQL option in the database engine.
func WithLogShowSQL(logShowSQL bool) EngineOpt {
	return func(e *engine) error {
		// set the log show SQL option for the database engine
		e.config.LogShowSQL = logShowSQL

		return nil
	}
}

// WithSkipCreation sets the skip creation logic in the database engine.
func WithSkipCreation(skipCreation bool) EngineOpt {
	return func(e *engine) error {
		// set to skip creating tables and indexes in the database engine
		e.config.SkipCreation = skipCreation

		return nil
	}
}

// WithTracing sets the shared tracing config in the database engine.
func WithTracing(tracing *tracing.Client) EngineOpt {
	return func(e *engine) error {
		e.tracing = tracing

		return nil
	}
}

// WithContext sets the context in the database engine.
func WithContext(ctx context.Context) EngineOpt {
	return func(e *engine) error {
		e.ctx = ctx

		return nil
	}
}

// WithLogPartitioned sets the log partitioned flag in the database engine.
func WithLogPartitioned(partitioned bool) EngineOpt {
	return func(e *engine) error {
		// set whether log table is partitioned in the database engine
		e.config.LogPartitioned = partitioned

		return nil
	}
}

// WithLogPartitionPattern sets the log partition pattern in the database engine.
func WithLogPartitionPattern(pattern string) EngineOpt {
	return func(e *engine) error {
		// set the naming pattern for log table partitions in the database engine
		e.config.LogPartitionPattern = pattern

		return nil
	}
}

// WithLogPartitionSchema sets the log partition schema in the database engine.
func WithLogPartitionSchema(schema string) EngineOpt {
	return func(e *engine) error {
		// set the database schema containing log table partitions in the database engine
		e.config.LogPartitionSchema = schema

		return nil
	}
}
