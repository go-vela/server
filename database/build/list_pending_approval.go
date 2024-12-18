// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListPendingAndRunningBuilds gets a list of all pending and running builds in the provided timeframe from the database.
func (e *engine) ListPendingApprovalBuilds(ctx context.Context, before string) ([]*api.Build, error) {
	e.logger.Trace("listing all pending and running builds")

	// variables to store query results and return value
	b := new([]types.Build)
	builds := []*api.Build{}

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Preload("Repo").
		Table(constants.TableBuild).
		Where("builds.created < ?", before).
		Where("builds.status = 'pending approval'").
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
			e.logger.Errorf("unable to decrypt repo: %v", err)
		}

		builds = append(builds, tmp.ToAPI())
	}

	return builds, nil
}
