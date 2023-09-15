// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

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
