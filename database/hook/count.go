// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package hook

import (
	"github.com/go-vela/types/constants"
)

// CountHooks gets the count of all hooks from the database.
func (e *engine) CountHooks() (int64, error) {
	e.logger.Tracef("getting count of all hooks from the database")

	// variable to store query results
	var h int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableHook).
		Count(&h).
		Error

	return h, err
}
