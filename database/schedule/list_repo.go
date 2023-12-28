// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// ListSchedulesForRepo gets a list of schedules by repo ID from the database.
func (e *engine) ListSchedulesForRepo(ctx context.Context, r *library.Repo, page, perPage int) ([]*library.Schedule, int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("listing schedules for repo %s from the database", r.GetFullName())

	// variables to store query results and return value
	count := int64(0)
	s := new([]database.Schedule)
	schedules := []*library.Schedule{}

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

		// convert query result to library type
		schedules = append(schedules, tmp.ToLibrary())
	}

	return schedules, count, nil
}
