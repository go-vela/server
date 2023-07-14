// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// ListActiveSchedules gets a list of all active schedules from the database.
func (e *engine) ListActiveSchedules() ([]*library.Schedule, error) {
	e.logger.Trace("listing all active schedules from the database")

	// variables to store query results and return value
	count := int64(0)
	s := new([]database.Schedule)
	schedules := []*library.Schedule{}

	// count the results
	count, err := e.CountActiveSchedules()
	if err != nil {
		return nil, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return schedules, nil
	}

	// send query to the database and store result in variable
	err = e.client.
		Table(constants.TableSchedule).
		Where("active = ?", true).
		Find(&s).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, schedule := range *s {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := schedule

		// convert query result to API type
		schedules = append(schedules, tmp.ToLibrary())
	}

	return schedules, nil
}
