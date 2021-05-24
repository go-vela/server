// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"github.com/go-vela/server/database/sqlite/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// GetRepoHookCount gets the count of webhooks by repo ID from the database.
func (c *client) GetRepoHookCount(r *library.Repo) (int64, error) {
	logrus.Tracef("getting count of hooks for repo %s from the database", r.GetFullName())

	// variable to store query results
	var h int64

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableHook).
		Raw(dml.SelectRepoHookCount, r.GetID()).
		Pluck("count", &h).Error

	return h, err
}
