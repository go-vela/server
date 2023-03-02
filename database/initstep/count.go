// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package initstep

import (
	"github.com/go-vela/types/constants"
)

// CountInitSteps gets the count of all inits from the database.
func (e *engine) CountInitSteps() (int64, error) {
	e.logger.Tracef("getting count of all init steps from the database")

	// variable to store query results
	var i int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableInitStep).
		Count(&i).
		Error

	return i, err
}
