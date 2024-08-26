// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// DeleteRepo deletes an existing repo from the database.
func (e *engine) DeleteRepo(ctx context.Context, r *api.Repo) error {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("deleting repo %s", r.GetFullName())

	// cast the library type to database type
	repo := types.RepoFromAPI(r)

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableRepo).
		Delete(repo).
		Error
}
