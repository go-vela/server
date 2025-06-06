// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-vela/server/constants"
)

type (
	// config represents the settings required to create the engine that implements the DeploymentInterface interface.
	config struct {
		// specifies the encryption key to use for the Hook engine
		EncryptionKey string
		// specifies to skip creating tables and indexes for the Deployment engine
		SkipCreation bool
	}

	// Engine represents the deployment functionality that implements the DeploymentInterface interface.
	Engine struct {
		// engine configuration settings used in deployment functions
		config *config

		ctx context.Context

		// gorm.io/gorm database client used in deployment functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in deployment functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}
)

// New creates and returns a Vela service for integrating with deployments in the database.
func New(opts ...EngineOpt) (*Engine, error) {
	// create new Deployment engine
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

	// check if we should skip creating deployment database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of deployment table and indexes")

		return e, nil
	}

	// create the deployments table
	err := e.CreateDeploymentTable(e.ctx, e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableDeployment, err)
	}

	// create the indexes for the deployments table
	err = e.CreateDeploymentIndexes(e.ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create indexes for %s table: %w", constants.TableDeployment, err)
	}

	return e, nil
}
