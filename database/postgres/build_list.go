// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"github.com/go-vela/server/database/postgres/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// GetBuildList gets a list of all builds from the database.
func (c *client) GetBuildList() ([]*library.Build, error) {
	logrus.Trace("listing builds from the database")

	// variable to store query results
	b := new([]database.Build)

	// send query to the database and store result in variable
	err := c.Postgres.
		Table(constants.TableBuild).
		Raw(dml.ListBuilds).
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

// GetOrgBuildList gets a list of all builds by org name from the database.
// nolint: lll // ignore long line length due to variable names
func (c *client) GetOrgBuildList(org string, exclude []int64, page, perPage int) ([]*library.Build, int64, error) {
	logrus.Tracef("listing builds for org %s from the database", org)

	// variable to store query results
	b := new([]database.Build)
	builds := []*library.Build{}
	count := int64(0)

	// // count the results
	count, err := c.GetOrgBuildCount(org, exclude)

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
	err = c.Postgres.
		Table(constants.TableBuild).
		Raw(dml.ListOrgBuilds, org, exclude, perPage, offset).
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
func (c *client) GetOrgBuildListByEvent(org string, exclude []int64, event string, page, perPage int) ([]*library.Build, int64, error) {
	logrus.Tracef("listing builds for org %s by event %s from the database", org, event)

	// variables to store query results
	b := new([]database.Build)
	builds := []*library.Build{}
	count := int64(0)

	// count the results
	count, err := c.GetOrgBuildCountByEvent(org, exclude, event)
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
	err = c.Postgres.
		Table(constants.TableBuild).
		Raw(dml.ListOrgBuildsByEvent, org, event, perPage, offset).
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

// GetRepoBuildList gets a list of all builds by repo ID from the database.
//
// nolint: lll // ignore long line length due to variable names
func (c *client) GetRepoBuildList(r *library.Repo, page, perPage int) ([]*library.Build, int64, error) {
	logrus.Tracef("listing builds for repo %s from the database", r.GetFullName())

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
	err = c.Postgres.
		Table(constants.TableBuild).
		Raw(dml.ListRepoBuilds, r.GetID(), perPage, offset).
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
	logrus.Tracef("listing builds for repo %s by event %s from the database", r.GetFullName(), event)

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
	err = c.Postgres.
		Table(constants.TableBuild).
		Raw(dml.ListRepoBuildsByEvent, r.GetID(), event, perPage, offset).
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
