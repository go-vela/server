// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package init

import (
	"github.com/go-vela/types/constants"
)

// CountInits gets the count of all inits from the database.
func (e *engine) CountInits() (int64, error) {
	e.logger.Tracef("getting count of all inits from the database")

	// variable to store query results
	var i int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableInit).
		Count(&i).
		Error

	return i, err
}
