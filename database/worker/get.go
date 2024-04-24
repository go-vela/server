// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// GetWorker gets a worker by ID from the database.
func (e *engine) GetWorker(ctx context.Context, id int64) (*api.Worker, error) {
	e.logger.Tracef("getting worker %d from the database", id)

	// variable to store query results
	w := new(types.Worker)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableWorker).
		Where("id = ?", id).
		Take(w).
		Error
	if err != nil {
		return nil, err
	}

	// return the worker
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Worker.ToLibrary
	return w.ToAPI(convertToBuilds(w.RunningBuildIDs)), nil
}
