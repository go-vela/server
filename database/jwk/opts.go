// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// EngineOpt represents a configuration option to initialize the database engine for key sets.
type EngineOpt func(*Engine) error

// WithClient sets the gorm.io/gorm client in the database engine for key sets.
func WithClient(client *gorm.DB) EngineOpt {
	return func(e *Engine) error {
		// set the gorm.io/gorm client in the key set engine
		e.client = client

		return nil
	}
}

// WithLogger sets the github.com/sirupsen/logrus logger in the database engine for key sets.
func WithLogger(logger *logrus.Entry) EngineOpt {
	return func(e *Engine) error {
		// set the github.com/sirupsen/logrus logger in the key set engine
		e.logger = logger

		return nil
	}
}

// WithSkipCreation sets the skip creation logic in the database engine for key sets.
func WithSkipCreation(skipCreation bool) EngineOpt {
	return func(e *Engine) error {
		// set to skip creating tables and indexes in the key set engine
		e.config.SkipCreation = skipCreation

		return nil
	}
}

// WithContext sets the context in the database engine for key sets.
func WithContext(ctx context.Context) EngineOpt {
	return func(e *Engine) error {
		e.ctx = ctx

		return nil
	}
}
