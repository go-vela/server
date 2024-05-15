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

// GetSchedule gets a schedule by ID from the database.
func (e *engine) GetSchedule(ctx context.Context, id int64) (*api.Schedule, error) {
	e.logger.Tracef("getting schedule %d from the database", id)

	// variable to store query results
	s := new(types.Schedule)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableSchedule).
		Preload("Repo").
		Preload("Repo.Owner").
		Where("id = ?", id).
		Take(s).
		Error
	if err != nil {
		return nil, err
	}

	// decrypt hash value for repo
	err = s.Repo.Decrypt(e.config.EncryptionKey)
	if err != nil {
		e.logger.Errorf("unable to decrypt repo %d: %v", s.Repo.ID.Int64, err)
	}

	// set repo to provided repo if creation successful
	result := s.ToAPI()

	// set next scheduled run
	t := time.Now().UTC()
	nextTime, _ := gronx.NextTickAfter(*result.Entry, t, false)
	result.SetNextRun(nextTime.Unix())

	return result, nil
}
