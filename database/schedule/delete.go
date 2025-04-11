// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// DeleteSchedule deletes an existing schedule from the database.
func (e *Engine) DeleteSchedule(ctx context.Context, s *api.Schedule) error {
	e.logger.WithFields(logrus.Fields{
		"schedule": s.GetName(),
	}).Tracef("deleting schedule %s in the database", s.GetName())

	// cast the API type to database type
	schedule := types.ScheduleFromAPI(s)

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableSchedule).
		Delete(schedule).
		Error
}
