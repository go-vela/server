// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// ListWorkers gets a list of all workers from the database.
func (e *engine) ListWorkers(ctx context.Context) ([]*api.Worker, error) {
	e.logger.Trace("listing all workers from the database")

	// variables to store query results and return value
	count := int64(0)
	w := new([]types.Worker)
	workers := []*api.Worker{}

	// count the results
	count, err := e.CountWorkers(ctx)
	if err != nil {
		return nil, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return workers, nil
	}

	// send query to the database and store result in variable
	err = e.client.
		Table(constants.TableWorker).
		Find(&w).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, worker := range *w {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := worker

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Worker.ToLibrary
		workers = append(workers, tmp.ToAPI(convertToBuilds(tmp.RunningBuildIDs)))
	}

	return workers, nil
}
