// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package executable

import (
	"context"

	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

// EngineOpt represents a configuration option to initialize the database engine for build executables.
type EngineOpt func(*engine) error

// WithClient sets the gorm.io/gorm client in the database engine for build executables.
func WithClient(client *gorm.DB) EngineOpt {
	return func(e *engine) error {
		// set the gorm.io/gorm client in the build executable engine
		e.client = client

		return nil
	}
}

// WithCompressionLevel sets the compression level in the database engine for build executables.
func WithCompressionLevel(level int) EngineOpt {
	return func(e *engine) error {
		// set the compression level in the build executable engine
		e.config.CompressionLevel = level

		return nil
	}
}

// WithEncryptionKey sets the encryption key in the database engine for build executables.
func WithEncryptionKey(key string) EngineOpt {
	return func(e *engine) error {
		// set the encryption key in the build executables engine
		e.config.EncryptionKey = key

		return nil
	}
}

// WithDriver sets the driver type in the database engine for build executables.
func WithDriver(driver string) EngineOpt {
	return func(e *engine) error {
		// set the driver type in the build executable engine
		e.config.Driver = driver

		return nil
	}
}

// WithLogger sets the github.com/sirupsen/logrus logger in the database engine for build executables.
func WithLogger(logger *logrus.Entry) EngineOpt {
	return func(e *engine) error {
		// set the github.com/sirupsen/logrus logger in the build executable engine
		e.logger = logger

		return nil
	}
}

// WithSkipCreation sets the skip creation logic in the database engine for build executables.
func WithSkipCreation(skipCreation bool) EngineOpt {
	return func(e *engine) error {
		// set to skip creating tables and indexes in the build executable engine
		e.config.SkipCreation = skipCreation

		return nil
	}
}

// WithContext sets the context in the database engine for build executables.
func WithContext(ctx context.Context) EngineOpt {
	return func(e *engine) error {
		e.ctx = ctx

		return nil
	}
}
