// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CountReposForUser gets the count of repos by user ID from the database.
func (e *engine) CountReposForUser(ctx context.Context, u *library.User, filters map[string]interface{}) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Tracef("getting count of repos for user %s from the database", u.GetName())

	// variable to store query results
	var r int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableRepo).
		Where("user_id = ?", u.GetID()).
		Where(filters).
		Count(&r).
		Error

	return r, err
}
