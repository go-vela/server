// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// EngineOpt represents a configuration option to initialize the database engine for Hooks.
type EngineOpt func(*Engine) error

// WithClient sets the gorm.io/gorm client in the database engine for Hooks.
func WithClient(client *gorm.DB) EngineOpt {
	return func(e *Engine) error {
		// set the gorm.io/gorm client in the hook engine
		e.client = client

		return nil
	}
}

// WithLogger sets the github.com/sirupsen/logrus logger in the database engine for Hooks.
func WithLogger(logger *logrus.Entry) EngineOpt {
	return func(e *Engine) error {
		// set the github.com/sirupsen/logrus logger in the hook engine
		e.logger = logger

		return nil
	}
}

// WithEncryptionKey sets the encryption key in the database engine for Builds.
func WithEncryptionKey(key string) EngineOpt {
	return func(e *Engine) error {
		// set the encryption key in the build engine
		e.config.EncryptionKey = key

		return nil
	}
}

// WithSkipCreation sets the skip creation logic in the database engine for Hooks.
func WithSkipCreation(skipCreation bool) EngineOpt {
	return func(e *Engine) error {
		// set to skip creating tables and indexes in the hook engine
		e.config.SkipCreation = skipCreation

		return nil
	}
}

// WithContext sets the context in the database engine for Repos.
func WithContext(ctx context.Context) EngineOpt {
	return func(e *Engine) error {
		e.ctx = ctx

		return nil
	}
}
