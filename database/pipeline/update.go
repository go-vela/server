// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// UpdatePipeline updates an existing pipeline in the database.
func (e *engine) UpdatePipeline(ctx context.Context, p *api.Pipeline) (*api.Pipeline, error) {
	e.logger.WithFields(logrus.Fields{
		"pipeline": p.GetCommit(),
	}).Tracef("updating pipeline %s in the database", p.GetCommit())

	pipeline := types.PipelineFromAPI(p)

	err := pipeline.Validate()
	if err != nil {
		return nil, err
	}

	err = pipeline.Compress(e.config.CompressionLevel)
	if err != nil {
		return nil, err
	}

	// send query to the database
	err = e.client.
		WithContext(ctx).
		Table(constants.TablePipeline).
		Save(pipeline).Error
	if err != nil {
		return nil, err
	}

	// decompress pipeline to return
	err = pipeline.Decompress()
	if err != nil {
		return nil, err
	}

	result := pipeline.ToAPI()
	result.SetRepo(p.GetRepo())

	return result, nil
}
