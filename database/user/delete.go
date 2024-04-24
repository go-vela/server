// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// DeleteUser deletes an existing user from the database.
func (e *engine) DeleteUser(ctx context.Context, u *api.User) error {
	e.logger.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Tracef("deleting user %s from the database", u.GetName())

	// cast the API type to database type
	user := types.UserFromAPI(u)

	// send query to the database
	return e.client.
		Table(constants.TableUser).
		Delete(user).
		Error
}
