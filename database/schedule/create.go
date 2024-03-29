// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code with update.go
package schedule

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CreateSchedule creates a new schedule in the database.
func (e *engine) CreateSchedule(ctx context.Context, s *library.Schedule) (*library.Schedule, error) {
	e.logger.WithFields(logrus.Fields{
		"schedule": s.GetName(),
	}).Tracef("creating schedule %s in the database", s.GetName())

	// cast the library type to database type
	schedule := database.ScheduleFromLibrary(s)

	// validate the necessary fields are populated
	err := schedule.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	result := e.client.Table(constants.TableSchedule).Create(schedule)

	return schedule.ToLibrary(), result.Error
}
