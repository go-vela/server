// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package init

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CountInitsForBuild gets the count of inits by build ID from the database.
func (e *engine) CountInitsForBuild(b *library.Build) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("getting count of inits for build %d from the database", b.GetNumber())

	// variable to store query results
	var i int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableInit).
		Where("build_id = ?", b.GetID()).
		Count(&i).
		Error

	return i, err
}
