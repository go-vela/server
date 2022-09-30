// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type (
	// config represents the settings required to create the engine that implements the RepoService interface.
	config struct {
		// specifies the encryption key to use for the Repo engine
		EncryptionKey string
		// specifies to skip creating tables and indexes for the Repo engine
		SkipCreation bool
	}

	// engine represents the repo functionality that implements the RepoService interface.
	engine struct {
		// engine configuration settings used in repo functions
		config *config

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
//
//nolint:revive // ignore returning unexported engine
func New(opts ...EngineOpt) (*engine, error) {
	// create new Repo engine
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

	// check if we should skip creating repo database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of repos table and indexes in the database")

		return e, nil
	}

	// create the repos table
	err := e.CreateRepoTable(e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableRepo, err)
	}

	// create the indexes for the repos table
	err = e.CreateRepoIndexes()
	if err != nil {
		return nil, fmt.Errorf("unable to create indexes for %s table: %w", constants.TableRepo, err)
	}

	return e, nil
}
