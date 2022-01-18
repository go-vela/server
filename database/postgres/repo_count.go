// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"github.com/go-vela/server/database/postgres/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// GetRepoCount gets a count of all repos from the database.
func (c *client) GetRepoCount() (int64, error) {
	c.Logger.Trace("getting count of repos from the database")

	// variable to store query results
	var r int64

	// send query to the database and store result in variable
	err := c.Postgres.
		Table(constants.TableRepo).
		Raw(dml.SelectReposCount).
		Pluck("count", &r).Error

	return r, err
}

// GetOrgRepoCount gets a count of all repos for a specific org from the database.
func (c *client) GetOrgRepoCount(org string, filters map[string]string) (int64, error) {
	c.Logger.WithFields(logrus.Fields{
		"org": org,
	}).Tracef("getting count of repos for org %s from the database", org)

	// variable to store query results
	var r int64

	// send query to the database and store result in variable
	err := c.Postgres.
		Table(constants.TableRepo).
		Select("count(*)").
		Where("org = ?", org).
		Where(filters).
		Pluck("count", &r).Error

	return r, err
}

// GetUserRepoCount gets a count of all repos for a specific user from the database.
func (c *client) GetUserRepoCount(u *library.User) (int64, error) {
	c.Logger.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Tracef("getting count of repos for user %s in the database", u.GetName())

	// variable to store query results
	var r int64

	// send query to the database and store result in variable
	err := c.Postgres.
		Table(constants.TableRepo).
		Raw(dml.SelectUserReposCount, u.GetID()).
		Pluck("count", &r).Error

	return r, err
}
