// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// CountBuildsForRepo gets the count of builds by repo ID from the database.
func (e *Engine) CountBuildsForRepo(ctx context.Context, r *api.Repo, filters map[string]interface{}, before, after int64) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("getting count of builds for repo %s", r.GetFullName())

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableBuild).
		Where("repo_id = ?", r.GetID()).
		Where("created < ?", before).
		Where("created > ?", after).
		Where(filters).
		Count(&b).
		Error

	return b, err
}
