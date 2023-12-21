// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"
	"strconv"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// ListWorkers gets a list of all workers from the database.
func (e *engine) ListWorkers(ctx context.Context, active string, before, after int64) ([]*library.Worker, error) {
	e.logger.Trace("listing all workers from the database")

	// variables to store query results and return value
	w := new([]database.Worker)
	workers := []*library.Worker{}

	// build query with checked in constraints
	query := e.client.Table(constants.TableWorker).
		Where("last_checked_in < ?", before).
		Where("last_checked_in > ?", after)

	// if active can be parsed as a boolean, add to query
	if _, err := strconv.ParseBool(active); err == nil {
		query.Where("active = ?", active)
	}

	// send query to the database and store result in variable
	err := query.Find(&w).Error
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
		workers = append(workers, tmp.ToLibrary())
	}

	return workers, nil
}
