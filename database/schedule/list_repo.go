// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/constants"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// ListSchedulesForRepo gets a list of schedules by repo ID from the database.
func (e *engine) ListSchedulesForRepo(r *library.Repo, page, perPage int) ([]*api.Schedule, int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("listing schedules for repo %s from the database", r.GetFullName())

	// variables to store query results and return value
	count := int64(0)
	h := new([]types.Schedule)
	schedules := []*api.Schedule{}

	// count the results
	count, err := e.CountSchedulesForRepo(r)
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
		Find(&h).
		Error
	if err != nil {
		return nil, count, err
	}

	// iterate through all query results
	for _, schedule := range *h {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := schedule

		// convert query result to library type
		schedules = append(schedules, tmp.ToAPI(r))
	}

	return schedules, count, nil
}
