// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// GetBuild gets a build by unique ID and repo ID from the database.
func (c *client) GetBuild(number int, r *library.Repo) (*library.Build, error) {
	logrus.Tracef("Getting build %s/%d from the database", r.GetFullName(), number)

	// variable to store query results
	b := new(database.Build)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableBuild).
		Raw(c.DML.BuildService.Select["repo"], r.GetID(), number).
		Scan(b).Error

	return b.ToLibrary(), err
}

// GetLastBuild gets the last build ran by repo ID from the database.
func (c *client) GetLastBuild(r *library.Repo) (*library.Build, error) {
	logrus.Tracef("Getting last build for repo %s from the database", r.GetFullName())

	// variable to store query results
	b := new(database.Build)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableBuild).
		Raw(c.DML.BuildService.Select["last"], r.GetID()).
		Scan(b).Error

	// the record will not exist if it's a new repo
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return b.ToLibrary(), err
}

// GetLastBuildByBranch gets the last build ran by repo ID and branch from the database.
func (c *client) GetLastBuildByBranch(r *library.Repo, branch string) (*library.Build, error) {
	logrus.Tracef("Getting last build for repo %s from the database", r.GetFullName())

	// variable to store query results
	b := new(database.Build)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableBuild).
		Raw(c.DML.BuildService.Select["lastByBranch"], r.GetID(), branch).
		Scan(b).Error

	// the record will not exist if it's a new repo
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return b.ToLibrary(), err
}

// GetBuildCount gets the count of all builds from the database.
func (c *client) GetBuildCount() (int64, error) {
	logrus.Trace("Count of builds from the database")

	// variable to store query results
	var b []int64

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableBuild).
		Raw(c.DML.BuildService.Select["count"]).
		Pluck("count", &b).Error

	return b[0], err
}

// GetBuildCountByStatus gets the count of all builds in a status from the database.
func (c *client) GetBuildCountByStatus(s string) (int64, error) {
	logrus.Trace("Count of builds by status from the database")

	// variable to store query results
	var b []int64

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableBuild).
		Raw(c.DML.BuildService.Select["countByStatus"], s).
		Pluck("count", &b).Error

	return b[0], err
}

// GetBuildList gets a list of all builds from the database.
func (c *client) GetBuildList() ([]*library.Build, error) {
	logrus.Trace("Listing builds from the database")

	// variable to store query results
	b := new([]database.Build)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableBuild).
		Raw(c.DML.BuildService.List["all"]).
		Scan(b).Error

	// variable we want to return
	builds := []*library.Build{}
	// iterate through all query results
	for _, build := range *b {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := build

		// convert query result to library type
		builds = append(builds, tmp.ToLibrary())
	}

	return builds, err
}

// GetRepoBuildList gets a list of all builds by repo ID from the database.
//
// nolint: lll // ignore long line length due to variable names
func (c *client) GetRepoBuildList(r *library.Repo, page, perPage int) ([]*library.Build, int64, error) {
	logrus.Tracef("Listing builds for repo %s from the database", r.GetFullName())

	// variable to store query results
	b := new([]database.Build)
	builds := []*library.Build{}
	count := int64(0)

	// count the results
	count, err := c.GetRepoBuildCount(r)
	if err != nil {
		return builds, 0, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return builds, 0, nil
	}

	// calculate offset for pagination through results
	offset := (perPage * (page - 1))

	// send query to the database and store result in variable
	err = c.Database.
		Table(constants.TableBuild).
		Raw(c.DML.BuildService.List["repo"], r.GetID(), perPage, offset).
		Scan(b).Error

	// iterate through all query results
	for _, build := range *b {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := build

		// convert query result to library type
		builds = append(builds, tmp.ToLibrary())
	}

	return builds, count, err
}

// GetOrgBuildList gets a list of all builds by org name from the database.
func (c *client) GetOrgBuildList(o string, page, perPage int) ([]*library.Build, int64, error) {
	logrus.Tracef("Listing builds for org %s from the database", o)

	// variable to store query results
	b := new([]database.Build)
	builds := []*library.Build{}
	count := int64(0)

	// // count the results
	count, err := c.GetOrgBuildCount(o)

	if err != nil {
		return builds, 0, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return builds, 0, nil
	}

	// calculate offset for pagination through results
	offset := (perPage * (page - 1))

	// send query to the database and store result in variable
	err = c.Database.
		Table(constants.TableBuild).
		Raw(c.DML.BuildService.List["org"], o, perPage, offset).
		Scan(b).Error

	// iterate through all query results
	for _, build := range *b {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := build

		// convert query result to library type
		builds = append(builds, tmp.ToLibrary())
	}

	return builds, count, err
}

// GetRepoBuildListByEvent gets a list of all builds by repo ID and event type from the database.
//
// nolint: lll // ignore long line length due to variable names
func (c *client) GetRepoBuildListByEvent(r *library.Repo, event string, page, perPage int) ([]*library.Build, int64, error) {
	logrus.Tracef("Listing builds for repo %s from the database by event '%s'", r.GetFullName(), event)

	// variables to store query results
	b := new([]database.Build)
	builds := []*library.Build{}
	count := int64(0)

	// count the results
	count, err := c.GetRepoBuildCountByEvent(r, event)
	if err != nil {
		return builds, 0, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return builds, 0, nil
	}

	// calculate offset for pagination through results
	offset := (perPage * (page - 1))

	// send query to the database and store result in variable
	err = c.Database.
		Table(constants.TableBuild).
		Raw(c.DML.BuildService.List["repoByEvent"], r.GetID(), event, perPage, offset).
		Scan(b).Error

	// iterate through all query results
	for _, build := range *b {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := build

		// convert query result to library type
		builds = append(builds, tmp.ToLibrary())
	}

	return builds, count, err
}

// GetOrgBuildListByEvent gets a list of all builds by org name and event type from the database.
//
// nolint: lll // ignore long line length due to variable names
func (c *client) GetOrgBuildListByEvent(org, event string, page, perPage int) ([]*library.Build, int64, error) {
	logrus.Tracef("Listing builds for repo %s from the database by event '%s'", org, event)

	// variables to store query results
	b := new([]database.Build)
	builds := []*library.Build{}
	count := int64(0)

	// count the results
	count, err := c.GetOrgBuildCountByEvent(org, event)
	if err != nil {
		return builds, 0, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return builds, 0, nil
	}

	// calculate offset for pagination through results
	offset := (perPage * (page - 1))

	// send query to the database and store result in variable
	err = c.Database.
		Table(constants.TableBuild).
		Raw(c.DML.BuildService.List["orgByEvent"], org, event, perPage, offset).
		Scan(b).Error

	// iterate through all query results
	for _, build := range *b {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := build

		// convert query result to library type
		builds = append(builds, tmp.ToLibrary())
	}

	return builds, count, err
}

// GetRepoBuildCount gets the count of all builds by repo ID from the database.
func (c *client) GetRepoBuildCount(r *library.Repo) (int64, error) {
	logrus.Trace("Count of builds from the database")

	// variable to store query results
	var b []int64

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableBuild).
		Raw(c.DML.BuildService.Select["countByRepo"], r.GetID()).
		Pluck("count", &b).Error

	return b[0], err
}

// GetOrgBuildCount gets the count of all builds by repo ID from the database.
func (c *client) GetOrgBuildCount(org string) (int64, error) {
	logrus.Trace("Count of builds from the database")

	// variable to store query results
	var b []int64

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableBuild).
		Raw(c.DML.BuildService.Select["countByOrg"], org).
		Pluck("count", &b).Error

	return b[0], err
}

// GetRepoBuildCountByEvent gets the count of all builds by repo ID and event from the database.
func (c *client) GetRepoBuildCountByEvent(r *library.Repo, event string) (int64, error) {
	logrus.Trace("Count of builds from the database")

	// variable to store query results
	var b []int64

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableBuild).
		Raw(c.DML.BuildService.Select["countByRepoAndEvent"], r.GetID(), event).
		Pluck("count", &b).Error

	return b[0], err
}

// GetOrgBuildCountByEvent gets the count of all builds by org name and event from the database.
func (c *client) GetOrgBuildCountByEvent(org string, event string) (int64, error) {
	logrus.Trace("Count of builds from the database")

	// variable to store query results
	var b []int64

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableBuild).
		Raw(c.DML.BuildService.Select["countByOrgAndEvent"], org, event).
		Pluck("count", &b).Error

	return b[0], err
}

// CreateBuild creates a new build in the database.
func (c *client) CreateBuild(b *library.Build) error {
	logrus.Tracef("Creating build %d in the database", b.GetNumber())

	// cast to database type
	build := database.BuildFromLibrary(b)

	// validate the necessary fields are populated
	err := build.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Database.
		Table(constants.TableBuild).
		Create(build.Crop()).Error
}

// UpdateBuild updates a build in the database.
func (c *client) UpdateBuild(b *library.Build) error {
	logrus.Tracef("Updating build %d in the database", b.GetNumber())

	// cast to database type
	build := database.BuildFromLibrary(b)

	// validate the necessary fields are populated
	err := build.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Database.
		Table(constants.TableBuild).
		Where("id = ?", b.GetID()).
		Update(build.Crop()).Error
}

// DeleteBuild deletes a build by unique ID from the database.
func (c *client) DeleteBuild(id int64) error {
	logrus.Tracef("Deleting build %d in the database", id)

	// send query to the database
	return c.Database.
		Table(constants.TableBuild).
		Exec(c.DML.BuildService.Delete, id).Error
}

// GetPendingAndRunningBuilds returns the list of pending and running builds
// within the given timeframe.
func (c *client) GetPendingAndRunningBuilds(after string) ([]*library.BuildQueue, error) {
	logrus.Trace("Selecting pending and running builds")

	// variable to store query results
	b := new([]database.BuildQueue)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableBuild).
		Raw(c.DML.BuildService.Select["pendingAndRunning"], after).
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
