// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// UpdatePipeline updates an existing pipeline in the database.
func (e *engine) UpdatePipeline(ctx context.Context, p *library.Pipeline) (*library.Pipeline, error) {
	e.logger.WithFields(logrus.Fields{
		"pipeline": p.GetCommit(),
	}).Tracef("updating pipeline %s in the database", p.GetCommit())

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
	err = e.client.Table(constants.TablePipeline).Save(pipeline).Error
	if err != nil {
		return nil, err
	}

	// decompress pipeline to return
	err = pipeline.Decompress()
	if err != nil {
		return nil, err
	}

	return pipeline.ToLibrary(), nil
}
