// SPDX-License-Identifier: Apache-2.0

package testreports

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// CountByRepo gets the count of all test reports by repo ID from the database.
func (e *Engine) CountByRepo(ctx context.Context, r *api.Repo, filters map[string]interface{}) (int64, error) {
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
		Joins("JOIN builds ON testreports.build_id = builds.id").
		Joins("JOIN repos ON builds.repo_id = repos.id").
		Where("repo_id = ?", r.GetID()).
		Where(filters).
		Count(&s).
		Error

	return s, err
}
