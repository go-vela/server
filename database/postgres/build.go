// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"errors"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database/postgres/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"gorm.io/gorm"
)

// GetBuild gets a build by number and repo ID from the database.
func (c *client) GetBuild(number int, r *library.Repo) (*library.Build, error) {
	c.Logger.WithFields(logrus.Fields{
		"build": number,
		"org":   r.GetOrg(),
		"repo":  r.GetName(),
	}).Tracef("getting build %s/%d from the database", r.GetFullName(), number)

	// variable to store query results
	b := new(database.Build)

	// send query to the database and store result in variable
	result := c.Postgres.
		Table(constants.TableBuild).
		Raw(dml.SelectRepoBuild, r.GetID(), number).
		Scan(b)

	// check if the query returned a record not found error or no rows were returned
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return b.ToLibrary(), result.Error
}

// GetBuildByID gets a build by id from the database.
func (c *client) GetBuildByID(id int64) (*library.Build, error) {
	c.Logger.WithFields(logrus.Fields{
		"build": id,
	}).Tracef("getting build %d from the database", id)

	// variable to store query result
	b := new(database.Build)

	// send query to the database and store result in variable
	result := c.Postgres.
		Table(constants.TableBuild).
		Raw(dml.SelectBuildByID, id).
		Scan(b)

	// check if the query returned a record not found error or no rows were returned
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return b.ToLibrary(), result.Error
}

// GetLastBuild gets the last build by repo ID from the database.
func (c *client) GetLastBuild(r *library.Repo) (*library.Build, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("getting last build for repo %s from the database", r.GetFullName())

	// variable to store query results
	b := new(database.Build)

	// send query to the database and store result in variable
	result := c.Postgres.
		Table(constants.TableBuild).
		Raw(dml.SelectLastRepoBuild, r.GetID()).
		Scan(b)

	// check if the query returned a record not found error or no rows were returned
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		// the record will not exist if it's a new repo
		return nil, nil
	}

	return b.ToLibrary(), result.Error
}

// GetLastBuildByBranch gets the last build by repo ID and branch from the database.
func (c *client) GetLastBuildByBranch(r *library.Repo, branch string) (*library.Build, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("getting last build for repo %s on branch %s from the database", r.GetFullName(), branch)

	// variable to store query results
	b := new(database.Build)

	// send query to the database and store result in variable
	result := c.Postgres.
		Table(constants.TableBuild).
		Raw(dml.SelectLastRepoBuildByBranch, r.GetID(), branch).
		Scan(b)

	// check if the query returned a record not found error or no rows were returned
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		// the record will not exist if it's a new repo
		return nil, nil
	}

	return b.ToLibrary(), result.Error
}

// GetPendingAndRunningBuilds returns the list of pending
// and running builds within the given timeframe.
func (c *client) GetPendingAndRunningBuilds(after string) ([]*library.BuildQueue, error) {
	c.Logger.Trace("getting pending and running builds from the database")

	// variable to store query results
	b := new([]database.BuildQueue)

	// send query to the database and store result in variable
	result := c.Postgres.
		Table(constants.TableBuild).
		Raw(dml.SelectPendingAndRunningBuilds, after).
		Scan(b)

	// variable we want to return
	builds := []*library.BuildQueue{}

	// iterate through all query results
	for _, build := range *b {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := build

		// convert query result to library type
		builds = append(builds, tmp.ToLibrary())
	}

	return builds, result.Error
}

// CreateBuild creates a new build in the database.
//
//nolint:dupl // ignore similar code to update.
func (c *client) CreateBuild(b *library.Build) (*library.Build, error) {
	c.Logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("creating build %d in the database", b.GetNumber())

	// cast to database type
	build := database.BuildFromLibrary(b)

	// validate the necessary fields are populated
	err := build.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	err = c.Postgres.
		Table(constants.TableBuild).
		Create(build.Crop()).Error

	if err != nil {
		return nil, err
	}

	return build.Crop().ToLibrary(), nil
}

// UpdateBuild updates a build in the database.
//
//nolint:dupl // ignore similar code with create.
func (c *client) UpdateBuild(b *library.Build) (*library.Build, error) {
	c.Logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("updating build %d in the database", b.GetNumber())

	// cast to database type
	build := database.BuildFromLibrary(b)

	// validate the necessary fields are populated
	err := build.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	err = c.Postgres.
		Table(constants.TableBuild).
		Save(build.Crop()).Error

	if err != nil {
		return nil, err
	}

	return build.Crop().ToLibrary(), nil
}

// DeleteBuild deletes a build by unique ID from the database.
func (c *client) DeleteBuild(id int64) error {
	c.Logger.Tracef("deleting build %d in the database", id)

	// send query to the database
	return c.Postgres.
		Table(constants.TableBuild).
		Exec(dml.DeleteBuild, id).Error
}
