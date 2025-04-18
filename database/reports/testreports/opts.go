// SPDX-License-Identifier: Apache-2.0

package testreports

import (
	"context"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// EngineOpt represents a configuration option to initialize the database Engine for TestReports.
type EngineOpt func(*Engine) error

// WithClient sets the gorm.io/gorm client in the database Engine for TestReports.
func WithClient(client *gorm.DB) EngineOpt {
	return func(e *Engine) error {
		// set the gorm.io/gorm client in the testreports Engine
		e.client = client

		return nil
	}
}

// WithEncryptionKey sets the encryption key in the database Engine for TestReports.
func WithEncryptionKey(key string) EngineOpt {
	return func(e *Engine) error {
		// set the encryption key in the testreports Engine
		e.config.EncryptionKey = key

		return nil
	}
}

// WithLogger sets the github.com/sirupsen/logrus logger in the database Engine for TestReports.
func WithLogger(logger *logrus.Entry) EngineOpt {
	return func(e *Engine) error {
		// set the github.com/sirupsen/logrus logger in the build Engine
		e.logger = logger

		return nil
	}
}

// WithSkipCreation sets the skip creation logic in the database Engine for TestReports.
func WithSkipCreation(skipCreation bool) EngineOpt {
	return func(e *Engine) error {
		// set to skip creating tables and indexes in the testreports Engine
		e.config.SkipCreation = skipCreation

		return nil
	}
}

// WithContext sets the context in the database Engine for TestReports.
func WithContext(ctx context.Context) EngineOpt {
	return func(e *Engine) error {
		e.ctx = ctx

		return nil
	}
}
