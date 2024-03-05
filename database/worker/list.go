// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"
	"fmt"
	"strconv"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
)

// ListWorkers gets a list of all workers from the database.
func (e *engine) ListWorkers(ctx context.Context, active string, before, after int64) ([]*api.Worker, error) {
	e.logger.Trace("listing all workers from the database")

	// variables to store query results and return value
	results := new([]Worker)
	workers := []*api.Worker{}

	// build query with checked in constraints
	query := e.client.Table(constants.TableWorker).
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
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := worker

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Worker.ToLibrary
		workers = append(workers, tmp.ToAPI(convertToBuilds(tmp.RunningBuildIDs)))
	}

	return workers, nil
}
