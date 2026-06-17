// SPDX-License-Identifier: Apache-2.0

package limits

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-vela/server/constants"
)

type (
	// config represents the settings required to create the engine that implements the LimitInterface interface.
	config struct {
		// specifies to skip creating tables and indexes for the Limit engine
		SkipCreation bool
	}

	// Engine represents the limit functionality that implements the LimitInterface interface.
	Engine struct {
		// engine configuration settings used in limit functions
		config *config

		ctx context.Context

		// gorm.io/gorm database client used in limit functions
		client *gorm.DB

		// sirupsen/logrus logger used in limit functions
		logger *logrus.Entry
	}
)

// New creates and returns a Vela service for integrating with limits in the database.
func New(opts ...EngineOpt) (*Engine, error) {
	e := new(Engine)

	e.client = new(gorm.DB)
	e.config = new(config)
	e.logger = new(logrus.Entry)

	for _, opt := range opts {
		err := opt(e)
		if err != nil {
			return nil, err
		}
	}

	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of org_build_limits table")

		return e, nil
	}

	err := e.CreateOrgBuildLimitTable(e.ctx, e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableOrgBuildLimit, err)
	}

	return e, nil
}
