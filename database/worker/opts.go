// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// EngineOpt represents a configuration option to initialize the database engine for Workers.
type EngineOpt func(*Engine) error

// WithClient sets the gorm.io/gorm client in the database engine for Workers.
func WithClient(client *gorm.DB) EngineOpt {
	return func(e *Engine) error {
		// set the gorm.io/gorm client in the worker engine
		e.client = client

		return nil
	}
}

// WithLogger sets the github.com/sirupsen/logrus logger in the database engine for Workers.
func WithLogger(logger *logrus.Entry) EngineOpt {
	return func(e *Engine) error {
		// set the github.com/sirupsen/logrus logger in the worker engine
		e.logger = logger

		return nil
	}
}

// WithSkipCreation sets the skip creation logic in the database engine for Workers.
func WithSkipCreation(skipCreation bool) EngineOpt {
	return func(e *Engine) error {
		// set to skip creating tables and indexes in the worker engine
		e.config.SkipCreation = skipCreation

		return nil
	}
}

// WithContext sets the context in the database engine for Workers.
func WithContext(ctx context.Context) EngineOpt {
	return func(e *Engine) error {
		e.ctx = ctx

		return nil
	}
}
