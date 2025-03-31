// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-vela/server/constants"
)

type (
	// config represents the settings required to create the engine that implements the JWKService interface.
	config struct {
		// specifies to skip creating tables and indexes for the JWK engine
		SkipCreation bool
		// specifies the driver for proper popping query
		Driver string
	}

	// engine represents the key set functionality that implements the JWKService interface.
	engine struct {
		// engine configuration settings used in key set functions
		config *config

		ctx context.Context

		// gorm.io/gorm database client used in key set functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in key set functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}
)

// New creates and returns a Vela service for integrating with key sets in the database.
//
//nolint:revive // ignore returning unexported engine
func New(opts ...EngineOpt) (*engine, error) {
	// create new JWK engine
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

	// check if we should skip creating key set database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of key sets table and indexes")

		return e, nil
	}

	// create the JWK table
	err := e.CreateJWKTable(e.ctx, e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableJWK, err)
	}

	return e, nil
}
