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

// GetScheduleForRepo gets a schedule by repo ID and number from the database.
func (e *engine) GetScheduleForRepo(r *library.Repo, name string) (*api.Schedule, error) {
	e.logger.WithFields(logrus.Fields{
		"org":      r.GetOrg(),
		"repo":     r.GetName(),
		"schedule": name,
	}).Tracef("getting schedule %s/%s from the database", r.GetFullName(), name)

	// variable to store query results
	s := new(types.Schedule)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableSchedule).
		Where("repo_id = ?", r.GetID()).
		Where("name = ?", name).
		Take(s).
		Error
	if err != nil {
		return nil, err
	}

	return s.ToAPI(r), nil
}
