// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
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
func (e *engine) CreatePipeline(p *library.Pipeline) (*library.Pipeline, error) {
	e.logger.WithFields(logrus.Fields{
		"pipeline": p.GetCommit(),
	}).Tracef("creating pipeline %s in the database", p.GetCommit())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#PipelineFromLibrary
	pipeline := database.PipelineFromLibrary(p)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Pipeline.Validate
	err := pipeline.Validate()
	if err != nil {
		return nil, err
	}

	// compress data for the pipeline
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Pipeline.Compress
	err = pipeline.Compress(e.config.CompressionLevel)
	if err != nil {
		return nil, err
	}

	// send query to the database
	err = e.client.Table(constants.TablePipeline).Create(pipeline).Error
	if err != nil {
		return nil, err
	}

	err = pipeline.Decompress()
	if err != nil {
		return nil, err
	}

	return pipeline.ToLibrary(), nil
}
