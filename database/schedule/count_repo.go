// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// CountSchedulesForRepo gets the count of schedules by repo ID from the database.
func (e *Engine) CountSchedulesForRepo(ctx context.Context, r *api.Repo) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("getting count of schedules for repo %s", r.GetFullName())

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableSchedule).
		Where("repo_id = ?", r.GetID()).
		Count(&s).
		Error

	return s, err
}
