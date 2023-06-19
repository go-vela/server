// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package user

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type (
	// config represents the settings required to create the engine that implements the UserInterface interface.
	config struct {
		// specifies the encryption key to use for the User engine
		EncryptionKey string
		// specifies to skip creating tables and indexes for the User engine
		SkipCreation bool
	}

	// engine represents the user functionality that implements the UserInterface interface.
	engine struct {
		// engine configuration settings used in user functions
		config *config

		// gorm.io/gorm database client used in user functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in user functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}
)

// New creates and returns a Vela service for integrating with users in the database.
//
//nolint:revive // ignore returning unexported engine
func New(opts ...EngineOpt) (UserInterface, error) {
	// create new User engine
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

	// check if we should skip creating user database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of users table and indexes in the database")

		return e, nil
	}

	// create the users table
	err := e.CreateUserTable(e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableUser, err)
	}

	// create the indexes for the users table
	err = e.CreateUserIndexes()
	if err != nil {
		return nil, fmt.Errorf("unable to create indexes for %s table: %w", constants.TableUser, err)
	}

	return e, nil
}
