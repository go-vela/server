// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"errors"

	"github.com/go-vela/server/database/sqlite/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

// GetBuild gets a build by number and repo ID from the database.
func (c *client) GetBuild(number int, r *library.Repo) (*library.Build, error) {
	logrus.Tracef("getting build %s/%d from the database", r.GetFullName(), number)

	// variable to store query results
	b := new(database.Build)

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableBuild).
		Raw(dml.SelectRepoBuild, r.GetID(), number).
		Scan(b).Error

	return b.ToLibrary(), err
}

// GetLastBuild gets the last build by repo ID from the database.
func (c *client) GetLastBuild(r *library.Repo) (*library.Build, error) {
	logrus.Tracef("getting last build for repo %s from the database", r.GetFullName())

	// variable to store query results
	b := new(database.Build)

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableBuild).
		Raw(dml.SelectLastRepoBuild, r.GetID()).
		Scan(b).Error

	// the record will not exist if it's a new repo
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return b.ToLibrary(), err
}

// GetLastBuildByBranch gets the last build by repo ID and branch from the database.
func (c *client) GetLastBuildByBranch(r *library.Repo, branch string) (*library.Build, error) {
	// nolint: lll // ignore long line length due to log message
	logrus.Tracef("getting last build for repo %s on branch %s from the database", r.GetFullName(), branch)

	// variable to store query results
	b := new(database.Build)

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableBuild).
		Raw(dml.SelectLastRepoBuildByBranch, r.GetID(), branch).
		Scan(b).Error

	// the record will not exist if it's a new repo
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return b.ToLibrary(), err
}

// GetPendingAndRunningBuilds returns the list of pending
// and running builds within the given timeframe.
func (c *client) GetPendingAndRunningBuilds(after string) ([]*library.BuildQueue, error) {
	logrus.Trace("getting pending and running builds from the database")

	// variable to store query results
	b := new([]database.BuildQueue)

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableBuild).
		Raw(dml.SelectPendingAndRunningBuilds, after).
		Scan(b).Error

	// variable we want to return
	builds := []*library.BuildQueue{}

	// iterate through all query results
	for _, build := range *b {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := build

		// convert query result to library type
		builds = append(builds, tmp.ToLibrary())
	}

	return builds, err
}

// CreateBuild creates a new build in the database.
func (c *client) CreateBuild(b *library.Build) error {
	logrus.Tracef("creating build %d in the database", b.GetNumber())

	// cast to database type
	build := database.BuildFromLibrary(b)

	// validate the necessary fields are populated
	err := build.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Sqlite.
		Table(constants.TableBuild).
		Create(build.Crop()).Error
}

// UpdateBuild updates a build in the database.
func (c *client) UpdateBuild(b *library.Build) error {
	logrus.Tracef("updating build %d in the database", b.GetNumber())

	// cast to database type
	build := database.BuildFromLibrary(b)

	// validate the necessary fields are populated
	err := build.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Sqlite.
		Table(constants.TableBuild).
		Save(build.Crop()).Error
}

// DeleteBuild deletes a build by unique ID from the database.
func (c *client) DeleteBuild(id int64) error {
	logrus.Tracef("deleting build %d in the database", id)

	// send query to the database
	return c.Sqlite.
		Table(constants.TableBuild).
		Exec(dml.DeleteBuild, id).Error
}
