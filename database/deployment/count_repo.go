// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package deployment

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CountDeploymentssForRepo gets the count of deployments by repo ID from the database.
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
