// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListUsers gets a list of all users from the database.
func (e *Engine) ListUsers(ctx context.Context) ([]*api.User, error) {
	e.logger.Trace("listing all users")

	// variables to store query results and return value
	u := new([]types.User)
	users := []*api.User{}

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableUser).
		Find(&u).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, user := range *u {
		// decrypt the fields for the user
		err = user.Decrypt(e.config.EncryptionKey)
		if err != nil {
			// TODO: remove backwards compatibility before 1.x.x release
			//
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch unencrypted users
			e.logger.Errorf("unable to decrypt user %d: %v", user.ID.Int64, err)
		}

		// convert query result to API type
		users = append(users, user.ToAPI())
	}

	return users, nil
}
