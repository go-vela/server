// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// CountBuildsForDeployment gets the count of builds by deployment URL from the database.
func (e *Engine) CountBuildsForDeployment(ctx context.Context, d *api.Deployment, filters map[string]interface{}) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"deployment": d.GetURL(),
	}).Tracef("getting count of builds for deployment %s", d.GetURL())

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableBuild).
		Where("source = ?", d.GetURL()).
		Where(filters).
		Order("number DESC").
		Count(&b).
		Error

	return b, err
}
