// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CountPipelinesForRepo gets the count of pipelines by repo ID from the database.
func (e *engine) CountPipelinesForRepo(ctx context.Context, r *library.Repo) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("getting count of pipelines for repo %s from the database", r.GetFullName())

	// variable to store query results
	var p int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TablePipeline).
		Where("repo_id = ?", r.GetID()).
		Count(&p).
		Error

	return p, err
}
