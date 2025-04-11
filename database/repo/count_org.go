// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
)

// CountReposForOrg gets the count of repos by org name from the database.
func (e *Engine) CountReposForOrg(ctx context.Context, org string, filters map[string]any) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org": org,
	}).Tracef("getting count of repos for org %s", org)

	// variable to store query results
	var r int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableRepo).
		Where("org = ?", org).
		Where(filters).
		Count(&r).
		Error

	return r, err
}
