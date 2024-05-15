// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"
	"time"

	"github.com/adhocore/gronx"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// ListSchedules gets a list of all schedules from the database.
func (e *engine) ListSchedules(ctx context.Context) ([]*api.Schedule, error) {
	e.logger.Trace("listing all schedules from the database")

	// variables to store query results and return value
	count := int64(0)
	s := new([]types.Schedule)
	schedules := []*api.Schedule{}

	// count the results
	count, err := e.CountSchedules(ctx)
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
		Preload("Repo").
		Preload("Repo.Owner").
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
			e.logger.Errorf("unable to decrypt repo %d: %v", tmp.ID.Int64, err)
		}

		// determine next scheduled run
		apiSchedule := tmp.ToAPI()
		t := time.Now().UTC()
		nextTime, _ := gronx.NextTickAfter(*apiSchedule.Entry, t, false)
		apiSchedule.SetNextRun(nextTime.Unix())

		// convert query result to API type
		schedules = append(schedules, apiSchedule)
	}

	return schedules, nil
}
