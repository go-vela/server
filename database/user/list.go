// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// ListUsers gets a list of all users from the database.
func (e *engine) ListUsers(ctx context.Context) ([]*library.User, error) {
	e.logger.Trace("listing all users from the database")

	// variables to store query results and return value
	count := int64(0)
	u := new([]database.User)
	users := []*library.User{}

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
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#User.Decrypt
		err = tmp.Decrypt(e.config.EncryptionKey)
		if err != nil {
			// TODO: remove backwards compatibility before 1.x.x release
			//
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch unencrypted users
			e.logger.Errorf("unable to decrypt user %d: %v", tmp.ID.Int64, err)
		}

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#User.ToLibrary
		users = append(users, tmp.ToLibrary())
	}

	return users, nil
}
