// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// EngineOpt represents a configuration option to initialize the database engine for Pipelines.
type EngineOpt func(*Engine) error

// WithClient sets the gorm.io/gorm client in the database engine for Pipelines.
func WithClient(client *gorm.DB) EngineOpt {
	return func(e *Engine) error {
		// set the gorm.io/gorm client in the pipeline engine
		e.client = client

		return nil
	}
}

// WithCompressionLevel sets the compression level in the database engine for Pipelines.
func WithCompressionLevel(level int) EngineOpt {
	return func(e *Engine) error {
		// set the compression level in the pipeline engine
		e.config.CompressionLevel = level

		return nil
	}
}

// WithEncryptionKey sets the encryption key in the database engine for Pipelines.
func WithEncryptionKey(key string) EngineOpt {
	return func(e *Engine) error {
		// set the encryption key in the build engine
		e.config.EncryptionKey = key

		return nil
	}
}

// WithLogger sets the github.com/sirupsen/logrus logger in the database engine for Pipelines.
func WithLogger(logger *logrus.Entry) EngineOpt {
	return func(e *Engine) error {
		// set the github.com/sirupsen/logrus logger in the pipeline engine
		e.logger = logger

		return nil
	}
}

// WithSkipCreation sets the skip creation logic in the database engine for Pipelines.
func WithSkipCreation(skipCreation bool) EngineOpt {
	return func(e *Engine) error {
		// set to skip creating tables and indexes in the pipeline engine
		e.config.SkipCreation = skipCreation

		return nil
	}
}

// WithContext sets the context in the database engine for Pipelines.
func WithContext(ctx context.Context) EngineOpt {
	return func(e *Engine) error {
		e.ctx = ctx

		return nil
	}
}
