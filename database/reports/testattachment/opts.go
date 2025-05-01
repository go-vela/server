// SPDX-License-Identifier: Apache-2.0

package testattachment

import (
	"context"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// EngineOpt represents a configuration option to initialize the database engine for test attachments.
type EngineOpt func(*Engine) error

// WithClient sets the gorm.io/gorm client in the database engine for test attachments.
func WithClient(client *gorm.DB) EngineOpt {
	return func(e *Engine) error {
		// set the gorm.io/gorm client in the testattachments engine
		e.client = client

		return nil
	}
}

// WithEncryptionKey sets the encryption key in the database engine for test attachments.
func WithEncryptionKey(key string) EngineOpt {
	return func(e *Engine) error {
		// set the encryption key in the testattachments engine
		e.config.EncryptionKey = key

		return nil
	}
}

// WithLogger sets the github.com/sirupsen/logrus logger in the database engine for test attachments.
func WithLogger(logger *logrus.Entry) EngineOpt {
	return func(e *Engine) error {
		// set the github.com/sirupsen/logrus logger in the testattachments engine
		e.logger = logger

		return nil
	}
}

// WithSkipCreation sets the skip creation logic in the database engine for test attachments.
func WithSkipCreation(skipCreation bool) EngineOpt {
	return func(e *Engine) error {
		// set to skip creating tables and indexes in the testattachments engine
		e.config.SkipCreation = skipCreation

		return nil
	}
}

// WithContext sets the context in the database engine for test attachments.
func WithContext(ctx context.Context) EngineOpt {
	return func(e *Engine) error {
		e.ctx = ctx

		return nil
	}
}
