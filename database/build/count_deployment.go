// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CountBuildsForDeployment gets the count of builds by deployment URL from the database.
func (e *engine) CountBuildsForDeployment(ctx context.Context, d *library.Deployment, filters map[string]interface{}) (int64, error) {
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
