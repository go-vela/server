// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
)

// CreateSchedule creates a new schedule in the database.
func (e *engine) CreateSchedule(ctx context.Context, s *api.Schedule) (*api.Schedule, error) {
	e.logger.WithFields(logrus.Fields{
		"schedule": s.GetName(),
	}).Tracef("creating schedule %s in the database", s.GetName())

	// cast the library type to database type
	schedule := FromAPI(s)

	// validate the necessary fields are populated
	err := schedule.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	result := e.client.Table(constants.TableSchedule).Create(schedule)

	return schedule.ToAPI(), result.Error
}
