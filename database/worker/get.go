// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetWorker gets a worker by ID from the database.
func (e *engine) GetWorker(id int64) (*library.Worker, error) {
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
