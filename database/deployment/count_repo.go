// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// CountDeploymentsForRepo gets the count of deployments by repo ID from the database.
func (e *Engine) CountDeploymentsForRepo(ctx context.Context, r *api.Repo) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("getting count of deployments for repo %s", r.GetFullName())

	// variable to store query results
	var d int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableDeployment).
		Where("repo_id = ?", r.GetID()).
		Count(&d).
		Error

	return d, err
}
