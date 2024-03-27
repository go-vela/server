// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// DeleteRepo deletes an existing repo from the database.
func (e *engine) DeleteRepo(ctx context.Context, r *api.Repo) error {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("deleting repo %s from the database", r.GetFullName())

	// cast the library type to database type
	repo := FromAPI(r)

	// send query to the database
	return e.client.
		Table(constants.TableRepo).
		Delete(repo).
		Error
}
