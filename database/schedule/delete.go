// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// DeleteSchedule deletes an existing schedule from the database.
func (e *engine) DeleteSchedule(ctx context.Context, s *library.Schedule) error {
	e.logger.WithFields(logrus.Fields{
		"schedule": s.GetName(),
	}).Tracef("deleting schedule %s in the database", s.GetName())

	// cast the library type to database type
	schedule := database.ScheduleFromLibrary(s)

	// send query to the database
	return e.client.
		Table(constants.TableSchedule).
		Delete(schedule).
		Error
}
