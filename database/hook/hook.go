// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-vela/server/constants"
)

type (
	// config represents the settings required to create the engine that implements the HookInterface interface.
	config struct {
		// specifies the encryption key to use for the Hook engine
		EncryptionKey string
		// specifies to skip creating tables and indexes for the Hook engine
		SkipCreation bool
	}

	// Engine represents the hook functionality that implements the HookInterface interface.
	Engine struct {
		// engine configuration settings used in hook functions
		config *config

		ctx context.Context

		// gorm.io/gorm database client used in hook functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in hook functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}
)

// New creates and returns a Vela service for integrating with hooks in the database.
func New(opts ...EngineOpt) (*Engine, error) {
	// create new Hook engine
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

	// check if we should skip creating hook database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of hooks table and indexes")

		return e, nil
	}

	// create the hooks table
	err := e.CreateHookTable(e.ctx, e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableHook, err)
	}

	// create the indexes for the hooks table
	err = e.CreateHookIndexes(e.ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create indexes for %s table: %w", constants.TableHook, err)
	}

	return e, nil
}
