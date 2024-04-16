// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
)

// DeleteSchedule deletes an existing schedule from the database.
func (e *engine) DeleteSchedule(ctx context.Context, s *api.Schedule) error {
	e.logger.WithFields(logrus.Fields{
		"schedule": s.GetName(),
	}).Tracef("deleting schedule %s in the database", s.GetName())

	// cast the library type to database type
	schedule := FromAPI(s)

	// send query to the database
	return e.client.
		Table(constants.TableSchedule).
		Delete(schedule).
		Error
}
