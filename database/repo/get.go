// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetRepo gets a repo by ID from the database.
func (e *Engine) GetRepo(ctx context.Context, id int64) (*api.Repo, error) {
	e.logger.Tracef("getting repo %d", id)

	// variable to store query results
	r := new(types.Repo)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableRepo).
		Preload("Owner").
		Where("id = ?", id).
		Take(r).
		Error
	if err != nil {
		return nil, err
	}

	// decrypt the fields for the repo
	err = r.Decrypt(e.config.EncryptionKey)
	if err != nil {
		// TODO: remove backwards compatibility before 1.x.x release
		//
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allows us to fetch unencrypted repos
		e.logger.Errorf("unable to decrypt repo %d: %v", id, err)

		// return the unencrypted repo
		return r.ToAPI(), nil
	}

	// return the decrypted repo
	return r.ToAPI(), nil
}
