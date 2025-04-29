// SPDX-License-Identifier: Apache-2.0

package testattachments

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-vela/server/constants"
)

type (
	// config represents the settings required to create the engine that implements the TestAttachmentsInterface interface.
	config struct {
		// specifies the encryption key to use for the TestAttachments engine
		EncryptionKey string
		// specifies to skip creating tables and indexes for the TestAttachments engine
		SkipCreation bool
	}

	// engine represents the testattachments functionality that implements the AttachmentInterface interface.
	Engine struct {
		// engine configuration settings used in testattachments functions
		config *config

		ctx context.Context

		// gorm.io/gorm database client used in testattachments functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in testattachments functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}
)

// New creates and returns a Vela service for integrating with testattachments in the database.
//

func New(opts ...EngineOpt) (*Engine, error) {
	// create new TestAttachments engine
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

	// check if we should skip creating testattachments database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of testattachments table and indexes")

		return e, nil
	}

	// create the testattachments table
	err := e.CreateTestAttachmentsTable(e.ctx, e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableAttachments, err)
	}

	// create the indexes for the testattachments table
	err = e.CreateTestAttachmentsIndexes(e.ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create indexes for %s table: %w", constants.TableAttachments, err)
	}

	return e, nil
}
