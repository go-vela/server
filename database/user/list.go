// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// ListUsers gets a list of all users from the database.
func (e *engine) ListUsers(ctx context.Context) ([]*api.User, error) {
	e.logger.Trace("listing all users")

	// variables to store query results and return value
	count := int64(0)
	u := new([]types.User)
	users := []*api.User{}

	// count the results
	count, err := e.CountUsers(ctx)
	if err != nil {
		return nil, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return users, nil
	}

	// send query to the database and store result in variable
	err = e.client.
		Table(constants.TableUser).
		Find(&u).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, user := range *u {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := user

		// decrypt the fields for the user
		err = tmp.Decrypt(e.config.EncryptionKey)
		if err != nil {
			// TODO: remove backwards compatibility before 1.x.x release
			//
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch unencrypted users
			e.logger.Errorf("unable to decrypt user %d: %v", tmp.ID.Int64, err)
		}

		// convert query result to API type
		users = append(users, tmp.ToAPI())
	}

	return users, nil
}
