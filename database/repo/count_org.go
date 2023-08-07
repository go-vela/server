// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// CountReposForOrg gets the count of repos by org name from the database.
func (e *engine) CountReposForOrg(org string, filters map[string]interface{}) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org": org,
	}).Tracef("getting count of repos for org %s from the database", org)

	// variable to store query results
	var r int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableRepo).
		Where("org = ?", org).
		Where(filters).
		Count(&r).
		Error

	return r, err
}
