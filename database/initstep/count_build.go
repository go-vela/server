// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package initstep

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CountInitStepsForBuild gets the count of inits by build ID from the database.
func (e *engine) CountInitStepsForBuild(b *library.Build) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"repo_id": b.GetRepoID(),
		"build":   b.GetNumber(),
	}).Tracef("getting count of init steps for build %d from the database", b.GetID())

	// variable to store query results
	var i int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableInitStep).
		Where("build_id = ?", b.GetID()).
		Count(&i).
		Error

	return i, err
}
