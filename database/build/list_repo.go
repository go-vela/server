// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// ListBuildsForRepo gets a list of builds by repo ID from the database.
//
//nolint:lll // ignore long line length due to variable names
func (e *engine) ListBuildsForRepo(ctx context.Context, r *library.Repo, filters map[string]interface{}, before, after int64, page, perPage int) ([]*library.Build, int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("listing builds for repo %s from the database", r.GetFullName())

	// variables to store query results and return values
	count := int64(0)
	b := new([]database.Build)
	builds := []*library.Build{}

	// count the results
	count, err := e.CountBuildsForRepo(ctx, r, filters)
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
		Where("repo_id = ?", r.GetID()).
		Where("created < ?", before).
		Where("created > ?", after).
		Where(filters).
		Order("number DESC").
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
