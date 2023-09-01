// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package user

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// DeleteUser deletes an existing user from the database.
func (e *engine) DeleteUser(ctx context.Context, u *library.User) error {
	e.logger.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Tracef("deleting user %s from the database", u.GetName())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#UserFromLibrary
	user := database.UserFromLibrary(u)

	// send query to the database
	return e.client.
		Table(constants.TableUser).
		Delete(user).
		Error
}
