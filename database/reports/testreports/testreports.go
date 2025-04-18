// SPDX-License-Identifier: Apache-2.0

package testreports

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-vela/server/constants"
)

type (
	// config represents the settings required to create the engine that implements the TestReportsInterface interface.
	config struct {
		// specifies the encryption key to use for the TestReports engine
		EncryptionKey string
		// specifies to skip creating tables and indexes for the TestReports engine
		SkipCreation bool
	}

	// engine represents the testreports functionality that implements the TestReportsInterface interface.
	Engine struct {
		// engine configuration settings used in testreports functions
		config *config

		ctx context.Context

		// gorm.io/gorm database client used in testreports functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in testreports functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}
)

// New creates and returns a Vela service for integrating with testreports in the database.
//
//nolint:revive // ignore returning unexported engine
func New(opts ...EngineOpt) (*Engine, error) {
	// create new TestReports engine
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

	// check if we should skip creating testreports database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of testreports table and indexes")

		return e, nil
	}

	// create the testreports table
	err := e.CreateTestReportsTable(e.ctx, e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableTestReports, err)
	}

	// create the indexes for the testreports table
	err = e.CreateTestReportsIndexes(e.ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create indexes for %s table: %w", constants.TableTestReports, err)
	}

	return e, nil
}
