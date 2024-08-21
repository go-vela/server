// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// ListActiveSchedules gets a list of all active schedules from the database.
func (e *engine) ListActiveSchedules(ctx context.Context) ([]*api.Schedule, error) {
	e.logger.Trace("listing all active schedules")

	// variables to store query results and return value
	count := int64(0)
	s := new([]types.Schedule)
	schedules := []*api.Schedule{}

	// count the results
	count, err := e.CountActiveSchedules(ctx)
	if err != nil {
		return nil, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return schedules, nil
	}

	// send query to the database and store result in variable
	err = e.client.
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
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := schedule

		// decrypt hash value for repo
		err = tmp.Repo.Decrypt(e.config.EncryptionKey)
		if err != nil {
			e.logger.Errorf("unable to decrypt repo %d: %v", tmp.Repo.ID.Int64, err)
		}

		// convert query result to API type
		schedules = append(schedules, tmp.ToAPI())
	}

	return schedules, nil
}
