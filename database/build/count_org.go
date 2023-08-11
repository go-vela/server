// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// CountBuildsForOrg gets the count of builds by org name from the database.
func (e *engine) CountBuildsForOrg(ctx context.Context, org string, filters map[string]interface{}) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org": org,
	}).Tracef("getting count of builds for org %s from the database", org)

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableBuild).
		Joins("JOIN repos ON builds.repo_id = repos.id").
		Where("repos.org = ?", org).
		Where(filters).
		Count(&b).
		Error

	return b, err
}
