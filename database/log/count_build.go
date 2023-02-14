// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package log

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

// CountLogsForBuild gets the count of logs by build ID from the database.
func (e *engine) CountLogsForBuild(b *library.Build) (int64, error) {
	e.logger.Tracef("getting count of logs for build %d from the database", b.GetID())

	// variable to store query results
	var l int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableLog).
		Where("build_id = ?", b.GetID()).
		Count(&l).
		Error

	return l, err
}
