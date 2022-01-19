// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

// engine represents the pipeline functionality that implements the PipelineService interface.
type engine struct {
	// specifies the level of compression to use for the Data field in a Pipeline.
	compressionLevel int

	// gorm.io database client used in pipeline functions
	//
	// https://pkg.go.dev/gorm.io/gorm#DB
	client *gorm.DB

	// sirupsen/logrus log entry used in pipeline functions
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
	logger *logrus.Entry
}

// New creates and returns a Vela service for integrating with pipelines in the database.
//
// nolint: revive // ignore returning unexported engine
func New(client *gorm.DB, logger *logrus.Entry, level int) *engine {
	return &engine{
		compressionLevel: level,
		client:           client,
		logger:           logger,
	}
}
