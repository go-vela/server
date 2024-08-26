// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/constants"
)

// CountBuildsForOrg gets the count of builds by org name from the database.
func (e *engine) CountBuildsForOrg(ctx context.Context, org string, filters map[string]interface{}) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org": org,
	}).Tracef("getting count of builds for org %s", org)

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableBuild).
		Joins("JOIN repos ON builds.repo_id = repos.id").
		Where("repos.org = ?", org).
		Where(filters).
		Count(&b).
		Error

	return b, err
}
