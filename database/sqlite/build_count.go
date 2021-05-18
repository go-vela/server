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

// GetBuildCount gets a count of all builds from the database.
func (c *client) GetBuildCount() (int64, error) {
	logrus.Trace("getting count of builds from the database")

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableBuild).
		Raw(dml.SelectBuildsCount).
		Pluck("count", &b).Error

	return b, err
}

// GetBuildCountByStatus gets a count of all builds by status from the database.
func (c *client) GetBuildCountByStatus(status string) (int64, error) {
	logrus.Tracef("getting count of builds by status %s from the database", status)

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableBuild).
		Raw(dml.SelectBuildsCountByStatus, status).
		Pluck("count", &b).Error

	return b, err
}

// GetOrgBuildCount gets the count of all builds by repo ID from the database.
func (c *client) GetOrgBuildCount(org string) (int64, error) {
	logrus.Tracef("getting count of builds for org %s from the database", org)

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableBuild).
		Raw(dml.SelectOrgBuildCount, org).
		Pluck("count", &b).Error

	return b, err
}

// GetOrgBuildCountByEvent gets the count of all builds by org name and event from the database.
func (c *client) GetOrgBuildCountByEvent(org string, event string) (int64, error) {
	logrus.Tracef("getting count of builds for org %s by event %s from the database", org, event)

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableBuild).
		Raw(dml.SelectOrgBuildCountByEvent, org, event).
		Pluck("count", &b).Error

	return b, err
}

// GetRepoBuildCount gets the count of all builds by repo ID from the database.
func (c *client) GetRepoBuildCount(r *library.Repo) (int64, error) {
	logrus.Tracef("getting count of builds for repo %s from the database", r.GetFullName())

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableBuild).
		Raw(dml.SelectRepoBuildCount, r.GetID()).
		Pluck("count", &b).Error

	return b, err
}

// GetRepoBuildCountByEvent gets the count of all builds by repo ID and event from the database.
func (c *client) GetRepoBuildCountByEvent(r *library.Repo, event string) (int64, error) {
	// nolint: lll // ignore long line length due to log message
	logrus.Tracef("getting count of builds for repo %s by event %s from the database", r.GetFullName(), event)

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableBuild).
		Raw(dml.SelectRepoBuildCountByEvent, r.GetID(), event).
		Pluck("count", &b).Error

	return b, err
}
