// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"github.com/go-vela/server/database/sqlite/dml"
	"github.com/go-vela/types/constants"
)

// GetWorkerCount gets a count of all workers from the database.
func (c *client) GetWorkerCount() (int64, error) {
	c.Logger.Trace("getting count of workers from the database")

	// variable to store query results
	var w int64

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableWorker).
		Raw(dml.SelectWorkersCount).
		Pluck("count", &w).Error

	return w, err
}
