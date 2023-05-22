// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CountBuildsForDeployment gets the count of builds by deployment URL from the database.
func (e *engine) CountBuildsForDeployment(d *library.Deployment, filters map[string]interface{}) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"deployment": d.GetURL(),
	}).Tracef("getting count of builds for deployment %s from the database", d.GetURL())

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableBuild).
		Where("source = ?", d.GetURL()).
		Where(filters).
		Order("number DESC").
		Count(&b).
		Error

	return b, err
}
