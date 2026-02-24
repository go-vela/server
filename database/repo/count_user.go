// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// CountReposForUser gets the count of repos by user ID from the database.
func (e *Engine) CountReposForUser(ctx context.Context, u *api.User, filters map[string]any) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Tracef("getting count of repos for user %s", u.GetName())

	// variable to store query results
	var r int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableRepo).
		Where("user_id = ?", u.GetID()).
		Where(filters).
		Count(&r).
		Error

	return r, err
}
