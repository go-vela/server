// SPDX-License-Identifier: Apache-2.0

package step

import (
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

// EngineOpt represents a configuration option to initialize the database engine for Steps.
type EngineOpt func(*engine) error

// WithClient sets the gorm.io/gorm client in the database engine for Steps.
func WithClient(client *gorm.DB) EngineOpt {
	return func(e *engine) error {
		// set the gorm.io/gorm client in the step engine
		e.client = client

		return nil
	}
}

// WithLogger sets the github.com/sirupsen/logrus logger in the database engine for Steps.
func WithLogger(logger *logrus.Entry) EngineOpt {
	return func(e *engine) error {
		// set the github.com/sirupsen/logrus logger in the step engine
		e.logger = logger

		return nil
	}
}

// WithSkipCreation sets the skip creation logic in the database engine for Steps.
func WithSkipCreation(skipCreation bool) EngineOpt {
	return func(e *engine) error {
		// set to skip creating tables and indexes in the step engine
		e.config.SkipCreation = skipCreation

		return nil
	}
}
