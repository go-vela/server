// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// ListBuildsForOrg gets a list of builds by org name from the database.
//
//nolint:lll // ignore long line length due to variable names
func (e *engine) ListBuildsForOrg(org string, filters map[string]interface{}, page, perPage int) ([]*library.Build, int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org": org,
	}).Tracef("listing builds for org %s from the database", org)

	// variables to store query results and return values
	count := int64(0)
	b := new([]database.Build)
	builds := []*library.Build{}

	// count the results
	count, err := e.CountBuildsForOrg(org, filters)
	if err != nil {
		return builds, 0, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return builds, 0, nil
	}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	err = e.client.
		Table(constants.TableBuild).
		Select("builds.*").
		Joins("JOIN repos ON builds.repo_id = repos.id").
		Where("repos.org = ?", org).
		Where(filters).
		Order("created DESC").
		Order("id").
		Limit(perPage).
		Offset(offset).
		Find(&b).
		Error
	if err != nil {
		return nil, count, err
	}

	// iterate through all query results
	for _, build := range *b {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := build

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Build.ToLibrary
		builds = append(builds, tmp.ToLibrary())
	}

	return builds, count, nil
}
