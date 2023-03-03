// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package initstep

import (
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

// EngineOpt represents a configuration option to initialize the database engine for InitSteps.
type EngineOpt func(*engine) error

// WithClient sets the gorm.io/gorm client in the database engine for InitSteps.
func WithClient(client *gorm.DB) EngineOpt {
	return func(e *engine) error {
		// set the gorm.io/gorm client in the init step engine
		e.client = client

		return nil
	}
}

// WithLogger sets the github.com/sirupsen/logrus logger in the database engine for InitSteps.
func WithLogger(logger *logrus.Entry) EngineOpt {
	return func(e *engine) error {
		// set the github.com/sirupsen/logrus logger in the init step engine
		e.logger = logger

		return nil
	}
}

// WithSkipCreation sets the skip creation logic in the database engine for InitSteps.
func WithSkipCreation(skipCreation bool) EngineOpt {
	return func(e *engine) error {
		// set to skip creating tables and indexes in the init step engine
		e.config.SkipCreation = skipCreation

		return nil
	}
}
