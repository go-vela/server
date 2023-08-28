// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

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
