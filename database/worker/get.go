// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetWorker gets a worker by ID from the database.
func (e *engine) GetWorker(ctx context.Context, id int64) (*library.Worker, error) {
	e.logger.Tracef("getting worker %d from the database", id)

	// variable to store query results
	w := new(database.Worker)

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
	return w.ToLibrary(), nil
}
