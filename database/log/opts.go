// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// EngineOpt represents a configuration option to initialize the database engine for Logs.
type EngineOpt func(*Engine) error

// WithClient sets the gorm.io/gorm client in the database engine for Logs.
func WithClient(client *gorm.DB) EngineOpt {
	return func(e *Engine) error {
		// set the gorm.io/gorm client in the log engine
		e.client = client

		return nil
	}
}

// WithCompressionLevel sets the compression level in the database engine for Logs.
func WithCompressionLevel(level int) EngineOpt {
	return func(e *Engine) error {
		// set the compression level in the log engine
		e.config.CompressionLevel = level

		return nil
	}
}

// WithLogger sets the github.com/sirupsen/logrus logger in the database engine for Logs.
func WithLogger(logger *logrus.Entry) EngineOpt {
	return func(e *Engine) error {
		// set the github.com/sirupsen/logrus logger in the log engine
		e.logger = logger

		return nil
	}
}

// WithSkipCreation sets the skip creation logic in the database engine for Logs.
func WithSkipCreation(skipCreation bool) EngineOpt {
	return func(e *Engine) error {
		// set to skip creating tables and indexes in the log engine
		e.config.SkipCreation = skipCreation

		return nil
	}
}

// WithContext sets the context in the database engine for Logs.
func WithContext(ctx context.Context) EngineOpt {
	return func(e *Engine) error {
		e.ctx = ctx

		return nil
	}
}

// WithLogPartitioned sets the log partitioned flag in the log engine.
func WithLogPartitioned(partitioned bool) EngineOpt {
	return func(e *Engine) error {
		e.config.LogPartitioned = partitioned

		return nil
	}
}

// WithLogPartitionPattern sets the log partition pattern in the log engine.
func WithLogPartitionPattern(pattern string) EngineOpt {
	return func(e *Engine) error {
		e.config.LogPartitionPattern = pattern

		return nil
	}
}

// WithLogPartitionSchema sets the log partition schema in the log engine.
func WithLogPartitionSchema(schema string) EngineOpt {
	return func(e *Engine) error {
		e.config.LogPartitionSchema = schema

		return nil
	}
}
