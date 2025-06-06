// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// EngineOpt represents a configuration option to initialize the database engine for dashboards.
type EngineOpt func(*Engine) error

// WithClient sets the gorm.io/gorm client in the database engine for dashboards.
func WithClient(client *gorm.DB) EngineOpt {
	return func(e *Engine) error {
		// set the gorm.io/gorm client in the dashboard engine
		e.client = client

		return nil
	}
}

// WithDriver sets the driver type in the database engine for dashboards.
func WithDriver(driver string) EngineOpt {
	return func(e *Engine) error {
		// set the driver type in the dashboard engine
		e.config.Driver = driver

		return nil
	}
}

// WithLogger sets the github.com/sirupsen/logrus logger in the database engine for dashboards.
func WithLogger(logger *logrus.Entry) EngineOpt {
	return func(e *Engine) error {
		// set the github.com/sirupsen/logrus logger in the dashboard engine
		e.logger = logger

		return nil
	}
}

// WithSkipCreation sets the skip creation logic in the database engine for dashboards.
func WithSkipCreation(skipCreation bool) EngineOpt {
	return func(e *Engine) error {
		// set to skip creating tables and indexes in the dashboard engine
		e.config.SkipCreation = skipCreation

		return nil
	}
}

// WithContext sets the context in the database engine for dashboards.
func WithContext(ctx context.Context) EngineOpt {
	return func(e *Engine) error {
		e.ctx = ctx

		return nil
	}
}
