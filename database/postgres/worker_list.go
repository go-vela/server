// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"github.com/go-vela/server/database/postgres/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetWorkerList gets a list of all workers from the database.
func (c *client) GetWorkerList() ([]*library.Worker, error) {
	c.Logger.Trace("listing workers from the database")

	// variable to store query results
	w := new([]database.Worker)

	// send query to the database and store result in variable
	err := c.Postgres.
		Table(constants.TableWorker).
		Raw(dml.ListWorkers).
		Scan(w).Error

	// variable we want to return
	workers := []*library.Worker{}
	// iterate through all query results
	for _, worker := range *w {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := worker

		// convert query result to library type
		workers = append(workers, tmp.ToLibrary())
	}

	return workers, err
}
