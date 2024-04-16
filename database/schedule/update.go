// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
	// "github.com/go-vela/types/database"
)

// UpdateSchedule updates an existing schedule in the database.
func (e *engine) UpdateSchedule(ctx context.Context, s *api.Schedule, fields bool) (*api.Schedule, error) {
	e.logger.WithFields(logrus.Fields{
		"schedule": s.GetName(),
	}).Tracef("updating schedule %s in the database", s.GetName())

	// cast the library type to database type
	schedule := FromAPI(s)

	// validate the necessary fields are populated
	err := schedule.Validate()
	if err != nil {
		return nil, err
	}

	// If "fields" is true, update entire record; otherwise, just update scheduled_at (part of processSchedule)
	//
	// we do this because Gorm will automatically set `updated_at` with the Save function
	// and the `updated_at` field should reflect the last time a user updated the record, rather than the scheduler
	if fields {
		err = e.client.Table(constants.TableSchedule).Save(schedule).Error
	} else {
		err = e.client.Table(constants.TableSchedule).Model(schedule).UpdateColumn("scheduled_at", s.GetScheduledAt()).Error
	}

	return schedule.ToAPI(), err
}
