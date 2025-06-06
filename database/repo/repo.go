// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-vela/server/constants"
)

type (
	// config represents the settings required to create the engine that implements the RepoInterface interface.
	config struct {
		// specifies the encryption key to use for the Repo engine
		EncryptionKey string
		// specifies to skip creating tables and indexes for the Repo engine
		SkipCreation bool
	}

	// Engine represents the repo functionality that implements the RepoInterface interface.
	Engine struct {
		// engine configuration settings used in repo functions
		config *config

		ctx context.Context

		// gorm.io/gorm database client used in repo functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in repo functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}
)

// New creates and returns a Vela service for integrating with repos in the database.
func New(opts ...EngineOpt) (*Engine, error) {
	// create new Repo engine
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

	// check if we should skip creating repo database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of repos table and indexes")

		return e, nil
	}

	// create the repos table
	err := e.CreateRepoTable(e.ctx, e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableRepo, err)
	}

	// create the indexes for the repos table
	err = e.CreateRepoIndexes(e.ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create indexes for %s table: %w", constants.TableRepo, err)
	}

	return e, nil
}
