// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
)

// CountHooksForRepo gets the count of hooks by repo ID from the database.
func (e *engine) CountHooksForRepo(ctx context.Context, r *api.Repo) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("getting count of hooks for repo %s", r.GetFullName())

	// variable to store query results
	var h int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableHook).
		Where("repo_id = ?", r.GetID()).
		Count(&h).
		Error

	return h, err
}
