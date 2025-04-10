// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListBuilds gets a list of all builds from the database.
func (e *Engine) ListBuilds(ctx context.Context) ([]*api.Build, error) {
	e.logger.Trace("listing all builds")

	// variables to store query results and return value
	b := new([]types.Build)
	builds := []*api.Build{}

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
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
		err = build.Repo.Decrypt(e.config.EncryptionKey)
		if err != nil {
			e.logger.Errorf("unable to decrypt repo: %v", err)
		}

		builds = append(builds, build.ToAPI())
	}

	return builds, nil
}
