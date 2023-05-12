// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"github.com/go-vela/types/constants"
)

// CountWorkersByStatus gets the count of all workers from the database with the specified status.
func (e *engine) CountWorkersByStatus(status string) (int64, error) {
	e.logger.Tracef("getting count of all workers from the database with the specified status")

	// variable to store query results
	var w int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableWorker).
		Where("status = ?", status).
		Count(&w).
		Error

	return w, err
}
