// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/constants"
	"github.com/go-vela/server/database/types"
	"github.com/sirupsen/logrus"
)

// DeleteSchedule deletes an existing schedule from the database.
func (e *engine) DeleteSchedule(s *api.Schedule) error {
	e.logger.WithFields(logrus.Fields{
		"org":      s.GetRepo().GetOrg(),
		"repo":     s.GetRepo().GetName(),
		"schedule": s.GetName(),
	}).Tracef("deleting schedule %s/%s in the database", s.GetRepo().GetFullName(), s.GetName())

	// cast the library type to database type
	schedule := types.ScheduleFromAPI(s)

	// send query to the database
	return e.client.
		Table(constants.TableSchedule).
		Delete(schedule).
		Error
}
