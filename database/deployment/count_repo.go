// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CountDeploymentsForRepo gets the count of deployments by repo ID from the database.
func (e *engine) CountDeploymentsForRepo(ctx context.Context, r *library.Repo) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("getting count of deployments for repo %s from the database", r.GetFullName())

	// variable to store query results
	var d int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableDeployment).
		Where("repo_id = ?", r.GetID()).
		Count(&d).
		Error

	return d, err
}
