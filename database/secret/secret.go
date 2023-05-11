// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type (
	// config represents the settings required to create the engine that implements the SecretService interface.
	config struct {
		// specifies the encryption key to use for the Secret engine
		EncryptionKey string
		// specifies to skip creating tables and indexes for the Secret engine
		SkipCreation bool
	}

	// engine represents the secret functionality that implements the SecretService interface.
	engine struct {
		// engine configuration settings used in secret functions
		config *config

		// gorm.io/gorm database client used in secret functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in secret functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}
)

// New creates and returns a Vela service for integrating with secrets in the database.
//
//nolint:revive // ignore returning unexported engine
func New(opts ...EngineOpt) (*engine, error) {
	// create new Secret engine
	e := new(engine)

	// create new fields
	e.client = new(gorm.DB)
	e.config = new(config)
	e.logger = new(logrus.Entry)

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(e)
		if err != nil {
			return nil, err
		}
	}

	// check if we should skip creating secret database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of secrets table and indexes in the database")

		return e, nil
	}

	// create the secrets table
	err := e.CreateSecretTable(e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableSecret, err)
	}

	// create the indexes for the secrets table
	err = e.CreateSecretIndexes()
	if err != nil {
		return nil, fmt.Errorf("unable to create indexes for %s table: %w", constants.TableSecret, err)
	}

	return e, nil
}
