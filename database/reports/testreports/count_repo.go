// SPDX-License-Identifier: Apache-2.0

package testreports

import (
	"context"
	api "github.com/go-vela/server/api/types"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
)

// CountByRepo gets the count of all test reports by repo ID from the database.
func (e *Engine) CountByRepo(ctx context.Context, r *api.Repo, filters map[string]interface{}, before, after int64) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("getting count of all test reports for repo %s", r.GetFullName())

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableTestReports).
		Where("repo_id = ?", r.GetID()).
		Where("created < ?", before).
		Where("created > ?", after).
		Where(filters).
		Count(&s).
		Error

	return s, err
}
