// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CreateSchedule creates a new schedule in the database.
func (e *Engine) CreateSchedule(ctx context.Context, s *api.Schedule) (*api.Schedule, error) {
	e.logger.WithFields(logrus.Fields{
		"schedule": s.GetName(),
	}).Tracef("creating schedule %s in the database", s.GetName())

	// cast the API type to database type
	schedule := types.ScheduleFromAPI(s)

	// validate the necessary fields are populated
	err := schedule.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	err = e.client.
		WithContext(ctx).
		Table(constants.TableSchedule).
		Create(schedule).Error
	if err != nil {
		return nil, err
	}

	// set repo to provided repo if creation successful
	result := schedule.ToAPI()
	result.SetRepo(s.GetRepo())

	return result, nil
}
