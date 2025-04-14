// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-vela/server/constants"
)

var (
	// ErrEmptyDashName defines the error type when a
	// User type has an empty Name field provided.
	ErrEmptyDashName = errors.New("empty dashboard name provided")

	// ErrExceededAdminLimit defines the error type when a
	// User type has Admins field provided that exceeds the database limit.
	ErrExceededAdminLimit = errors.New("exceeded admins limit")
)

type (
	// config represents the settings required to create the engine that implements the DashboardService interface.
	config struct {
		// specifies to skip creating tables and indexes for the Dashboard engine
		SkipCreation bool
		// specifies the driver for proper popping query
		Driver string
	}

	// Engine represents the dashboard functionality that implements the DashboardService interface.
	Engine struct {
		// engine configuration settings used in dashboard functions
		config *config

		ctx context.Context

		// gorm.io/gorm database client used in dashboard functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in dashboard functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}
)

// New creates and returns a Vela service for integrating with dashboards in the database.
func New(opts ...EngineOpt) (*Engine, error) {
	// create new Dashboard engine
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

	// check if we should skip creating dashboard database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of dashboards table and indexes")

		return e, nil
	}

	// create the dashboards table
	err := e.CreateDashboardTable(e.ctx, e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableDashboard, err)
	}

	return e, nil
}
