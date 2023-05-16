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
func (e *engine) ListBuildsForOrg(org, sortBy string, filters map[string]interface{}, page, perPage int) ([]*library.Build, int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org": org,
	}).Tracef("listing builds for org %s from the database", org)

	// variables to store query results and return values
	count := int64(0)
	r := new([]database.Build)
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

	switch sortBy {
	case "latest":
		query := e.client.
			Table(constants.TableBuild).
			Select("builds.id, MAX(builds.created) AS latest_build").
			Joins("INNER JOIN builds builds ON builds.repo_id = builds.id").
			Where("builds.org = ?", org).
			Group("builds.id")

		err = e.client.
			Table(constants.TableBuild).
			Select("builds.*").
			Joins("LEFT JOIN (?) t on builds.id = t.id", query).
			Order("latest_build DESC NULLS LAST").
			Limit(perPage).
			Offset(offset).
			Find(&r).
			Error
		if err != nil {
			return nil, count, err
		}
	case "name":
		fallthrough
	default:
		err = e.client.
			Table(constants.TableBuild).
			Where("org = ?", org).
			Where(filters).
			Order("name").
			Limit(perPage).
			Offset(offset).
			Find(&r).
			Error
		if err != nil {
			return nil, count, err
		}
	}

	// iterate through all query results
	for _, build := range *r {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := build

		// decrypt the fields for the build
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Build.Decrypt
		err = tmp.Decrypt(e.config.EncryptionKey)
		if err != nil {
			// TODO: remove backwards compatibility before 1.x.x release
			//
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch unencrypted builds
			e.logger.Errorf("unable to decrypt build %d: %v", tmp.ID.Int64, err)
		}

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Build.ToLibrary
		builds = append(builds, tmp.ToLibrary())
	}

	return builds, count, nil
}
