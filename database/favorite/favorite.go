// SPDX-License-Identifier: Apache-2.0

package favorite

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-vela/server/constants"
)

type (
	// config represents the settings required to create the engine that implements the UserInterface interface.
	config struct {
		// specifies to skip creating tables and indexes for the User engine
		SkipCreation bool
		Driver       string
	}

	// Engine represents the user functionality that implements the UserInterface interface.
	Engine struct {
		// engine configuration settings used in user functions
		config *config

		ctx context.Context

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
func New(opts ...EngineOpt) (*Engine, error) {
	// create new User engine
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

	// check if we should skip creating user database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of users table and indexes")

		return e, nil
	}

	// create the users table
	err := e.CreateFavoritesTable(e.ctx, e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableUser, err)
	}

	// create the indexes for the users table
	err = e.CreateFavoritesIndexes(e.ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create indexes for %s table: %w", constants.TableUser, err)
	}

	return e, nil
}
