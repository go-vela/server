// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListActiveSchedules gets a list of all active schedules from the database.
func (e *Engine) ListActiveSchedules(ctx context.Context) ([]*api.Schedule, error) {
	e.logger.Trace("listing all active schedules")

	// variables to store query results and return value
	s := new([]types.Schedule)
	schedules := []*api.Schedule{}

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableSchedule).
		Preload("Repo").
		Preload("Repo.Owner").
		Where("active = ?", true).
		Find(&s).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, schedule := range *s {
		// decrypt hash value for repo
		err = schedule.Repo.Decrypt(e.config.EncryptionKey)
		if err != nil {
			e.logger.Errorf("unable to decrypt repo %d: %v", schedule.Repo.ID.Int64, err)
		}

		// convert query result to API type
		schedules = append(schedules, schedule.ToAPI())
	}

	return schedules, nil
}
