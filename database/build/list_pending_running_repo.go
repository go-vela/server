// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// ListPendingAndRunningBuilds gets a list of all pending and running builds in the provided timeframe from the database.
func (e *engine) ListPendingAndRunningBuildsForRepo(ctx context.Context, repo *api.Repo) ([]*api.Build, error) {
	e.logger.Trace("listing all pending and running builds from the database")

	// variables to store query results and return value
	b := new([]types.Build)
	builds := []*api.Build{}

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableBuild).
		Preload("Repo").
		Preload("Repo.Owner").
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

		err = tmp.Repo.Decrypt(e.config.EncryptionKey)
		if err != nil {
			e.logger.Errorf("unable to decrypt repo %s/%s: %v", repo.GetOrg(), repo.GetName(), err)
		}

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Build.ToLibrary
		builds = append(builds, tmp.ToAPI())
	}

	return builds, nil
}
