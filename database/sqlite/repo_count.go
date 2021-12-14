// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"github.com/go-vela/server/database/sqlite/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

// GetRepoCount gets a count of all repos from the database.
func (c *client) GetRepoCount() (int64, error) {
	c.Logger.Trace("getting count of repos from the database")

	// variable to store query results
	var r int64

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableRepo).
		Raw(dml.SelectReposCount).
		Pluck("count", &r).Error

	return r, err
}

// GetOrgRepoCount gets a count of all repos for a specific org from the database.
func (c *client) GetOrgRepoCount(org string, filters map[string]string) (int64, error) {
	c.Logger.Tracef("getting count of repos for org %s in the database", org)

	// variable to store query results
	var r int64

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableRepo).
		Select("count(*)").
		Where("org = ?", org).
		Where(filters).
		Pluck("count", &r).Error

	return r, err
}

// GetUserRepoCount gets a count of all repos for a specific user from the database.
func (c *client) GetUserRepoCount(u *library.User) (int64, error) {
	c.Logger.Tracef("getting count of repos for user %s in the database", u.GetName())

	// variable to store query results
	var r int64

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableRepo).
		Raw(dml.SelectUserReposCount, u.GetID()).
		Pluck("count", &r).Error

	return r, err
}
