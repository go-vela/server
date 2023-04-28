// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/constants"
	"github.com/go-vela/server/database/types"
	"github.com/sirupsen/logrus"
)

// UpdateSchedule updates an existing schedule in the database.
func (e *engine) UpdateSchedule(s *api.Schedule) error {
	e.logger.WithFields(logrus.Fields{
		"org":      s.GetRepo().GetOrg(),
		"repo":     s.GetRepo().GetName(),
		"schedule": s.GetName(),
	}).Tracef("updating schedule %s/%s in the database", s.GetRepo().GetFullName(), s.GetName())

	// cast the library type to database type
	schedule := types.ScheduleFromAPI(s)

	// validate the necessary fields are populated
	err := schedule.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return e.client.
		Table(constants.TableSchedule).
		Save(schedule).
		Error
}
