// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// EngineOpt represents a configuration option to initialize the database engine for Schedules.
type EngineOpt func(*Engine) error

// WithClient sets the gorm.io/gorm client in the database engine for Schedules.
func WithClient(client *gorm.DB) EngineOpt {
	return func(e *Engine) error {
		// set the gorm.io/gorm client in the schedule engine
		e.client = client

		return nil
	}
}

// WithEncryptionKey sets the encryption key in the database engine for Schedules.
func WithEncryptionKey(key string) EngineOpt {
	return func(e *Engine) error {
		// set the encryption key in the schedule engine
		e.config.EncryptionKey = key

		return nil
	}
}

// WithLogger sets the github.com/sirupsen/logrus logger in the database engine for Schedules.
func WithLogger(logger *logrus.Entry) EngineOpt {
	return func(e *Engine) error {
		// set the github.com/sirupsen/logrus logger in the schedule engine
		e.logger = logger

		return nil
	}
}

// WithSkipCreation sets the skip creation logic in the database engine for Schedules.
func WithSkipCreation(skipCreation bool) EngineOpt {
	return func(e *Engine) error {
		// set to skip creating tables and indexes in the schedule engine
		e.config.SkipCreation = skipCreation

		return nil
	}
}

// WithContext sets the context in the database engine for Schedules.
func WithContext(ctx context.Context) EngineOpt {
	return func(e *Engine) error {
		e.ctx = ctx

		return nil
	}
}
