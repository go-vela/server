// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
)

// ListBuilds gets a list of all builds from the database.
func (e *engine) ListBuilds(ctx context.Context) ([]*api.Build, error) {
	e.logger.Trace("listing all builds from the database")

	// variables to store query results and return value
	count := int64(0)
	b := new([]Build)
	builds := []*api.Build{}

	// count the results
	count, err := e.CountBuilds(ctx)
	if err != nil {
		return nil, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return builds, nil
	}

	// send query to the database and store result in variable
	err = e.client.
		Preload("Repo").
		Preload("Repo.Owner").
		Table(constants.TableBuild).
		Find(&b).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, build := range *b {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := build

		builds = append(builds, tmp.ToAPI())
	}

	return builds, nil
}
