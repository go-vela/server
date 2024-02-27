// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// ListPendingAndRunningBuilds gets a list of all pending and running builds in the provided timeframe from the database.
func (e *engine) ListPendingAndRunningBuildsForRepo(ctx context.Context, repo *library.Repo) ([]*library.Build, error) {
	e.logger.Trace("listing all pending and running builds from the database")

	// variables to store query results and return value
	b := new([]database.Build)
	builds := []*library.Build{}

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableBuild).
		Select("*").
		Where("repo_id = ?", repo.GetID()).
		Where("status = 'running' OR status = 'pending' OR status = 'pending approval'").
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
