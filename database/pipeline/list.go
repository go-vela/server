// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListPipelines gets a list of all pipelines from the database.
func (e *engine) ListPipelines(ctx context.Context) ([]*api.Pipeline, error) {
	e.logger.Trace("listing all pipelines")

	// variables to store query results and return value
	count := int64(0)
	p := new([]types.Pipeline)
	pipelines := []*api.Pipeline{}

	// count the results
	count, err := e.CountPipelines(ctx)
	if err != nil {
		return nil, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return pipelines, nil
	}

	// send query to the database and store result in variable
	err = e.client.
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
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := pipeline

		err = tmp.Decompress()
		if err != nil {
			return nil, err
		}

		err = tmp.Repo.Decrypt(e.config.EncryptionKey)
		if err != nil {
			e.logger.Errorf("unable to decrypt repo: %v", err)
		}

		pipelines = append(pipelines, tmp.ToAPI())
	}

	return pipelines, nil
}
