// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// CountPipelinesForRepo gets the count of pipelines by repo ID from the database.
func (e *Engine) CountPipelinesForRepo(ctx context.Context, r *api.Repo) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("getting count of pipelines for repo %s", r.GetFullName())

	// variable to store query results
	var p int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TablePipeline).
		Where("repo_id = ?", r.GetID()).
		Count(&p).
		Error

	return p, err
}
