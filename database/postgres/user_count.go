// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"github.com/go-vela/server/database/postgres/dml"
	"github.com/go-vela/types/constants"
)

// GetUserCount gets a count of all users from the database.
func (c *client) GetUserCount() (int64, error) {
	c.Logger.Trace("getting count of users from the database")

	// variable to store query results
	var u int64

	// send query to the database and store result in variable
	err := c.Postgres.
		Table(constants.TableUser).
		Raw(dml.SelectUsersCount).
		Pluck("count", &u).Error

	return u, err
}
