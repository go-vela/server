// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"github.com/go-vela/types/constants"
)

// CountBuilds gets the count of all builds from the database.
func (e *engine) CountBuilds() (int64, error) {
	e.logger.Tracef("getting count of all builds from the database")

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableBuild).
		Count(&b).
		Error

	return b, err
}
