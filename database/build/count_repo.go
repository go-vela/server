// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CountBuildsForRepo gets the count of builds by repo ID from the database.
func (e *engine) CountBuildsForRepo(ctx context.Context, r *library.Repo, filters map[string]interface{}) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("getting count of builds for repo %s from the database", r.GetFullName())

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableBuild).
		Where("repo_id = ?", r.GetID()).
		Where(filters).
		Count(&b).
		Error

	return b, err
}
