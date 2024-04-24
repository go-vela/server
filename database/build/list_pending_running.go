// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// ListPendingAndRunningBuilds gets a list of all pending and running builds in the provided timeframe from the database.
func (e *engine) ListPendingAndRunningBuilds(ctx context.Context, after string) ([]*api.QueueBuild, error) {
	e.logger.Trace("listing all pending and running builds from the database")

	// variables to store query results and return value
	b := new([]types.QueueBuild)
	builds := []*api.QueueBuild{}

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableBuild).
		Select("builds.created, builds.number, builds.status, repos.full_name").
		InnerJoins("INNER JOIN repos ON builds.repo_id = repos.id").
		Where("builds.created > ?", after).
		Where("builds.status = 'running' OR builds.status = 'pending'").
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
