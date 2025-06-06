// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"

	"github.com/go-vela/server/constants"
)

// CountActiveSchedules gets the count of all active schedules from the database.
func (e *Engine) CountActiveSchedules(ctx context.Context) (int64, error) {
	e.logger.Tracef("getting count of all active schedules")

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableSchedule).
		Where("active = ?", true).
		Count(&s).
		Error

	return s, err
}
