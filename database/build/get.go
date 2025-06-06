// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetBuild gets a build by ID from the database.
func (e *Engine) GetBuild(ctx context.Context, id int64) (*api.Build, error) {
	e.logger.Tracef("getting build %d", id)

	// variable to store query results
	b := new(types.Build)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableBuild).
		Preload("Repo").
		Preload("Repo.Owner").
		Where("id = ?", id).
		Take(b).
		Error
	if err != nil {
		return nil, err
	}

	err = b.Repo.Decrypt(e.config.EncryptionKey)
	if err != nil {
		e.logger.Errorf("unable to decrypt repo: %v", err)
	}

	return b.ToAPI(), nil
}
