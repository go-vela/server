// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-vela/server/constants"
)

type (
	// config represents the settings required to create the engine that implements the SecretInterface interface.
	config struct {
		// specifies the encryption key to use for the Secret engine
		EncryptionKey string
		// specifies to skip creating tables and indexes for the Secret engine
		SkipCreation bool
	}

	// Engine represents the secret functionality that implements the SecretInterface interface.
	Engine struct {
		// engine configuration settings used in secret functions
		config *config

		ctx context.Context

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
func New(opts ...EngineOpt) (*Engine, error) {
	// create new Secret engine
	e := new(Engine)

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
		e.logger.Warning("skipping creation of secrets table and indexes")

		return e, nil
	}

	// create the secrets table
	err := e.CreateSecretTables(e.ctx, e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableSecret, err)
	}

	// create the indexes for the secrets table
	err = e.CreateSecretIndexes(e.ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create indexes for %s table: %w", constants.TableSecret, err)
	}

	return e, nil
}
