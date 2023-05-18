// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	"github.com/go-vela/types/constants"
)

// CountActiveSchedules gets the count of all active schedules from the database.
func (e *engine) CountActiveSchedules() (int64, error) {
	e.logger.Tracef("getting count of all active schedules from the database")

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableSchedule).
		Where("active = ?", true).
		Count(&s).
		Error

	return s, err
}
