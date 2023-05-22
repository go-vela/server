// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"github.com/go-vela/types/constants"
)

// CountBuildsForStatus gets the count of builds by org name from the database.
func (e *engine) CountBuildsForStatus(status string, filters map[string]interface{}) (int64, error) {
	e.logger.Tracef("getting count of builds for status %s from the database", status)

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableBuild).
		Where("status = ?", status).
		Where(filters).
		Count(&b).
		Error

	return b, err
}
