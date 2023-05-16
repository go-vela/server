// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CountBuildsForUser gets the count of builds by user ID from the database.
func (e *engine) CountBuildsForRepo(r *library.Repo, filters map[string]interface{}) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Tracef("getting count of builds for user %s from the database", u.GetName())

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableBuild).
		Where("user_id = ?", u.GetID()).
		Where(filters).
		Count(&r).
		Error

	return r, err
}
