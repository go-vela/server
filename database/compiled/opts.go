// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package compiled

import (
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

// EngineOpt represents a configuration option to initialize the database engine for Pipelines.
type EngineOpt func(*engine) error

// WithClient sets the gorm.io/gorm client in the database engine for Pipelines.
func WithClient(client *gorm.DB) EngineOpt {
	return func(e *engine) error {
		// set the gorm.io/gorm client in the pipeline engine
		e.client = client

		return nil
	}
}

// WithCompressionLevel sets the compression level in the database engine for Pipelines.
func WithCompressionLevel(level int) EngineOpt {
	return func(e *engine) error {
		// set the compression level in the pipeline engine
		e.config.CompressionLevel = level

		return nil
	}
}

// WithLogger sets the github.com/sirupsen/logrus logger in the database engine for Pipelines.
func WithLogger(logger *logrus.Entry) EngineOpt {
	return func(e *engine) error {
		// set the github.com/sirupsen/logrus logger in the pipeline engine
		e.logger = logger

		return nil
	}
}

// WithSkipCreation sets the skip creation logic in the database engine for Pipelines.
func WithSkipCreation(skipCreation bool) EngineOpt {
	return func(e *engine) error {
		// set to skip creating tables and indexes in the pipeline engine
		e.config.SkipCreation = skipCreation

		return nil
	}
}
