// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"

	"github.com/go-vela/server/constants"
)

// CountUsers gets the count of all users from the database.
func (e *Engine) CountUsers(ctx context.Context) (int64, error) {
	e.logger.Tracef("getting count of all users")

	// variable to store query results
	var u int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableUser).
		Count(&u).
		Error

	return u, err
}
