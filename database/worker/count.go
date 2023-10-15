// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"

	"github.com/go-vela/types/constants"
)

// CountWorkers gets the count of all workers from the database.
func (e *engine) CountWorkers(ctx context.Context) (int64, error) {
	e.logger.Tracef("getting count of all workers from the database")

	// variable to store query results
	var w int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableWorker).
		Count(&w).
		Error

	return w, err
}
