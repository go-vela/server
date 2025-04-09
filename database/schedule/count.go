// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"

	"github.com/go-vela/server/constants"
)

// CountSchedules gets the count of all schedules from the database.
func (e *Engine) CountSchedules(ctx context.Context) (int64, error) {
	e.logger.Tracef("getting count of all schedules")

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableSchedule).
		Count(&s).
		Error

	return s, err
}
