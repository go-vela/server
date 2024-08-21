// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// GetUser gets a user by ID from the database.
func (e *engine) GetUser(ctx context.Context, id int64) (*api.User, error) {
	e.logger.Tracef("getting user %d", id)

	// variable to store query results
	u := new(types.User)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableUser).
		Where("id = ?", id).
		Take(u).
		Error
	if err != nil {
		return nil, err
	}

	// decrypt the fields for the user
	err = u.Decrypt(e.config.EncryptionKey)
	if err != nil {
		// TODO: remove backwards compatibility before 1.x.x release
		//
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allows us to fetch unencrypted users
		e.logger.Errorf("unable to decrypt user %d: %v", u.ID.Int64, err)
	}

	// return the decrypted user
	return u.ToAPI(), nil
}
