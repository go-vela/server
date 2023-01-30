// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

// EngineOpt represents a configuration option to initialize the database engine for Repos.
type EngineOpt func(*engine) error

// WithClient sets the gorm.io/gorm client in the database engine for Repos.
func WithClient(client *gorm.DB) EngineOpt {
	return func(e *engine) error {
		// set the gorm.io/gorm client in the repo engine
		e.client = client

		return nil
	}
}

// WithEncryptionKey sets the encryption key in the database engine for Repos.
func WithEncryptionKey(key string) EngineOpt {
	return func(e *engine) error {
		// set the encryption key in the repo engine
		e.config.EncryptionKey = key

		return nil
	}
}

// WithLogger sets the github.com/sirupsen/logrus logger in the database engine for Repos.
func WithLogger(logger *logrus.Entry) EngineOpt {
	return func(e *engine) error {
		// set the github.com/sirupsen/logrus logger in the repo engine
		e.logger = logger

		return nil
	}
}

// WithSkipCreation sets the skip creation logic in the database engine for Repos.
func WithSkipCreation(skipCreation bool) EngineOpt {
	return func(e *engine) error {
		// set to skip creating tables and indexes in the repo engine
		e.config.SkipCreation = skipCreation

		return nil
	}
}
