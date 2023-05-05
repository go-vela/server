// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

//nolint:dupl // ignore similar code with update.go
package schedule

import (
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/constants"
	"github.com/go-vela/server/database/types"
	"github.com/sirupsen/logrus"
)

// CreateSchedule creates a new schedule in the database.
func (e *engine) CreateSchedule(s *api.Schedule) error {
	e.logger.WithFields(logrus.Fields{
		"org":      s.GetRepo().GetOrg(),
		"repo":     s.GetRepo().GetName(),
		"schedule": s.GetName(),
	}).Tracef("creating schedule %s/%s in the database", s.GetRepo().GetFullName(), s.GetName())

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
		Create(schedule).
		Error
}
