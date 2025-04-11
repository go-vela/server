// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"
	"fmt"
	"strconv"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListWorkers gets a list of all workers from the database.
func (e *Engine) ListWorkers(ctx context.Context, active string, before, after int64) ([]*api.Worker, error) {
	e.logger.Trace("listing all workers")

	// variables to store query results and return value
	results := new([]types.Worker)
	workers := []*api.Worker{}

	// build query with checked in constraints
	query := e.client.
		WithContext(ctx).
		Table(constants.TableWorker).
		Where("last_checked_in < ?", before).
		Where("last_checked_in > ?", after)

	// if active can be parsed as a boolean, add to query
	if b, err := strconv.ParseBool(active); err == nil {
		// convert bool to 0/1 for Sqlite
		qBool := 0
		if b {
			qBool = 1
		}

		query.Where("active = ?", fmt.Sprintf("%d", qBool))
	}

	// send query to the database and store result in variable
	err := query.Find(&results).Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, worker := range *results {
		// convert query result to API type
		workers = append(workers, worker.ToAPI(convertToBuilds(worker.RunningBuildIDs)))
	}

	return workers, nil
}
