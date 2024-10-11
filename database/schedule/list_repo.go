// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// ListSchedulesForRepo gets a list of schedules by repo ID from the database.
func (e *engine) ListSchedulesForRepo(ctx context.Context, r *api.Repo, page, perPage int) ([]*api.Schedule, int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("listing schedules for repo %s", r.GetFullName())

	// variables to store query results and return value
	count := int64(0)
	s := new([]types.Schedule)
	schedules := []*api.Schedule{}

	// count the results
	count, err := e.CountSchedulesForRepo(ctx, r)
	if err != nil {
		return nil, 0, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return schedules, 0, nil
	}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	// send query to the database and store result in variable
	err = e.client.
		WithContext(ctx).
		Table(constants.TableSchedule).
		Where("repo_id = ?", r.GetID()).
		Order("id DESC").
		Limit(perPage).
		Offset(offset).
		Find(&s).
		Error
	if err != nil {
		return nil, count, err
	}

	// iterate through all query results
	for _, schedule := range *s {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := schedule

		result := tmp.ToAPI()
		result.SetRepo(r)

		// convert query result to API type
		schedules = append(schedules, result)
	}

	return schedules, count, nil
}
