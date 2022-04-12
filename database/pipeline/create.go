// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CreatePipeline creates a new pipeline in the database.
func (e *engine) CreatePipeline(p *library.Pipeline) error {
	e.logger.WithFields(logrus.Fields{
		"pipeline": p.GetNumber(),
	}).Tracef("creating pipeline %d in the database", p.GetNumber())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#PipelineFromLibrary
	pipeline := database.PipelineFromLibrary(p)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Pipeline.Validate
	err := pipeline.Validate()
	if err != nil {
		return err
	}

	// compress data for the pipeline
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Pipeline.Compress
	err = pipeline.Compress(e.config.CompressionLevel)
	if err != nil {
		return err
	}

	// send query to the database
	return e.client.
		Table(constants.TablePipeline).
		Create(pipeline).
		Error
}
