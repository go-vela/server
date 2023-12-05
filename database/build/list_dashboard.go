// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// ListBuildsForDashboardRepo gets a list of builds by repo ID from the database.
//
//nolint:lll // ignore long line length due to variable names
func (e *engine) ListBuildsForDashboardRepo(ctx context.Context, r *library.Repo, branches, events []string) ([]*library.Build, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("listing builds for repo %s from the database", r.GetFullName())

	// variables to store query results and return values
	b := new([]database.Build)
	builds := []*library.Build{}

	query := e.client.Table(constants.TableBuild).Where("repo_id = ?", r.GetID())

	if len(branches) > 0 {
		query = query.Where("branch IN (?)", branches)
	}

	if len(events) > 0 {
		query = query.Where("event IN (?)", events)
	}

	err := query.
		Order("number DESC").
		Limit(5).
		Find(&b).
		Error
	if err != nil {
		return nil, err
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

	return builds, nil
}
