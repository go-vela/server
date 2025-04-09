// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// DeletePipeline deletes an existing pipeline from the database.
func (e *Engine) DeletePipeline(ctx context.Context, p *api.Pipeline) error {
	e.logger.WithFields(logrus.Fields{
		"pipeline": p.GetCommit(),
	}).Tracef("deleting pipeline %s", p.GetCommit())

	pipeline := types.PipelineFromAPI(p)

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TablePipeline).
		Delete(pipeline).
		Error
}
