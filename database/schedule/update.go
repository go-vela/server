// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// UpdateSchedule updates an existing schedule in the database.
func (e *engine) UpdateSchedule(s *library.Schedule, fields bool) error {
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

	// If "fields" is true, update entire record; otherwise, just update scheduled_at (part of processSchedule)
	//
	// we do this because Gorm will automatically set `updated_at` with the Save function
	// and the `updated_at` field should reflect the last time a user updated the record, rather than the scheduler
	if fields {
		err = e.client.Table(constants.TableSchedule).Save(schedule).Error
	} else {
		err = e.client.Table(constants.TableSchedule).Model(schedule).
			UpdateColumns(database.Schedule{ScheduledAt: schedule.ScheduledAt, Processing: schedule.Processing}).Error
	}

	return err
}
