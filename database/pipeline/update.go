// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// UpdatePipeline updates an existing pipeline in the database.
func (e *engine) UpdatePipeline(p *library.Pipeline) error {
	e.logger.WithFields(logrus.Fields{
		"pipeline": p.GetNumber(),
	}).Tracef("updating pipeline %d in the database", p.GetNumber())

	// cast to database type
	pipeline := database.PipelineFromLibrary(p)

	// validate the necessary fields are populated
	err := pipeline.Validate()
	if err != nil {
		return err
	}

	// compress data for the pipeline
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Log.Compress
	err = pipeline.Compress(e.compressionLevel)
	if err != nil {
		return err
	}

	// send query to the database
	return e.client.
		Table(constants.TablePipeline).
		Save(pipeline).
		Error
}
