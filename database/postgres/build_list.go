// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
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
	c.Logger.Trace("listing builds from the database")

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

// GetDeploymentBuildList gets a list of all builds from the database.
func (c *client) GetDeploymentBuildList(deployment string) ([]*library.Build, error) {
	c.Logger.WithFields(logrus.Fields{
		"deployment": deployment,
	}).Tracef("listing builds for deployment %s from the database", deployment)

	// variable to store query results
	b := new([]database.Build)

	filters := map[string]string{}
	if len(deployment) > 0 {
		filters["source"] = deployment
	}
	// send query to the database and store result in variable
	err := c.Postgres.
		Table(constants.TableBuild).
		Select("*").
		Where(filters).
		Limit(3).
		Order("number DESC").
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

// GetOrgBuildList gets a list of all builds by org name and allows filters from the database.
func (c *client) GetOrgBuildList(org string, filters map[string]interface{}, page, perPage int) ([]*library.Build, int64, error) {
	c.Logger.WithFields(logrus.Fields{
		"org": org,
	}).Tracef("listing builds for org %s from the database", org)

	// variables to store query results
	b := new([]database.Build)
	builds := []*library.Build{}
	count := int64(0)

	// count the results
	count, err := c.GetOrgBuildCount(org, filters)
	if err != nil {
		return builds, 0, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return builds, 0, nil
	}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	// send query to the database and store result in variable
	err = c.Postgres.
		Table(constants.TableBuild).
		Select("builds.*").
		Joins("JOIN repos ON builds.repo_id = repos.id and repos.org = ?", org).
		Where(filters).
		Order("created DESC").
		Order("id").
		Limit(perPage).
		Offset(offset).
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
func (c *client) GetRepoBuildList(r *library.Repo, filters map[string]interface{}, page, perPage int) ([]*library.Build, int64, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("listing builds for repo %s from the database", r.GetFullName())

	// variable to store query results
	b := new([]database.Build)
	builds := []*library.Build{}
	count := int64(0)

	// count the results
	count, err := c.GetRepoBuildCount(r, filters)
	if err != nil {
		return builds, 0, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return builds, 0, nil
	}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	err = c.Postgres.
		Table(constants.TableBuild).
		Where("repo_id = ?", r.GetID()).
		Where(filters).
		Order("number DESC").
		Limit(perPage).
		Offset(offset).
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
