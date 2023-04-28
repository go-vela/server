// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	"github.com/go-vela/server/database/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CountSchedulesForRepo gets the count of schedules by repo ID from the database.
func (e *engine) CountSchedulesForRepo(r *library.Repo) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("getting count of schedules for repo %s from the database", r.GetFullName())

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableSchedule).
		Where("repo_id = ?", r.GetID()).
		Count(&s).
		Error

	return s, err
}
