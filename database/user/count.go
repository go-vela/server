// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"

	"github.com/go-vela/types/constants"
)

// CountUsers gets the count of all users from the database.
func (e *engine) CountUsers(ctx context.Context) (int64, error) {
	e.logger.Tracef("getting count of all users from the database")

	// variable to store query results
	var u int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableUser).
		Count(&u).
		Error

	return u, err
}
