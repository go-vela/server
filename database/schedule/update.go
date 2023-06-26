// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

//nolint:dupl // ignore similar code with create.go
package schedule

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// UpdateSchedule updates an existing schedule in the database.
func (e *engine) UpdateSchedule(s *library.Schedule, config bool) error {
	e.logger.WithFields(logrus.Fields{
		"schedule": s.GetName(),
	}).Tracef("updating schedule %s in the database", s.GetName())

	// cast the library type to database type
	schedule := database.ScheduleFromLibrary(s)

	// validate the necessary fields are populated
	err := schedule.Validate()
	if err != nil {
		return err
	}

	// if update is just setting the scheduled_at, then ignore updating other fields
	if config {
		err = e.client.Table(constants.TableSchedule).Save(schedule).Error
	} else {
		err = e.client.Table(constants.TableSchedule).Model(schedule).UpdateColumn("scheduled_at", s.GetScheduledAt()).Error
	}

	return err
}
