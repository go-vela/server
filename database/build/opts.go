// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

// EngineOpt represents a configuration option to initialize the database engine for Builds.
type EngineOpt func(*engine) error

// WithClient sets the gorm.io/gorm client in the database engine for Builds.
func WithClient(client *gorm.DB) EngineOpt {
	return func(e *engine) error {
		// set the gorm.io/gorm client in the build engine
		e.client = client

		return nil
	}
}

// WithLogger sets the github.com/sirupsen/logrus logger in the database engine for Builds.
func WithLogger(logger *logrus.Entry) EngineOpt {
	return func(e *engine) error {
		// set the github.com/sirupsen/logrus logger in the build engine
		e.logger = logger

		return nil
	}
}

// WithSkipCreation sets the skip creation logic in the database engine for Builds.
func WithSkipCreation(skipCreation bool) EngineOpt {
	return func(e *engine) error {
		// set to skip creating tables and indexes in the build engine
		e.config.SkipCreation = skipCreation

		return nil
	}
}

// WithContext sets the context in the database engine for Builds.
func WithContext(ctx context.Context) EngineOpt {
	return func(e *engine) error {
		e.ctx = ctx

		return nil
	}
}
