// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListPipelines gets a list of all pipelines from the database.
func (e *Engine) ListPipelines(ctx context.Context) ([]*api.Pipeline, error) {
	e.logger.Trace("listing all pipelines")

	// variables to store query results and return value
	p := new([]types.Pipeline)
	pipelines := []*api.Pipeline{}

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TablePipeline).
		Preload("Repo").
		Preload("Repo.Owner").
		Find(&p).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, pipeline := range *p {
		err = pipeline.Decompress()
		if err != nil {
			return nil, err
		}

		err = pipeline.Repo.Decrypt(e.config.EncryptionKey)
		if err != nil {
			e.logger.Errorf("unable to decrypt repo: %v", err)
		}

		pipelines = append(pipelines, pipeline.ToAPI())
	}

	return pipelines, nil
}
