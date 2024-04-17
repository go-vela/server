// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
)

// GetSchedule gets a schedule by ID from the database.
func (e *engine) GetSchedule(ctx context.Context, id int64) (*api.Schedule, error) {
	e.logger.Tracef("getting schedule %d from the database", id)

	// variable to store query results
	s := new(Schedule)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableSchedule).
		Preload("Repo").
		Where("id = ?", id).
		Take(s).
		Error
	if err != nil {
		return nil, err
	}

	return s.ToAPI(), nil
}
