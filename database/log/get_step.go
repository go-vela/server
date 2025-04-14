// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetLogForStep gets a log by step ID from the database.
func (e *Engine) GetLogForStep(ctx context.Context, s *api.Step) (*api.Log, error) {
	e.logger.Tracef("getting log for step %d for build %d", s.GetID(), s.GetBuildID())

	// variable to store query results
	l := new(types.Log)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableLog).
		Where("step_id = ?", s.GetID()).
		Take(l).
		Error
	if err != nil {
		return nil, err
	}

	// decompress log data for the step
	err = l.Decompress()
	if err != nil {
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allows us to fetch uncompressed logs
		e.logger.Errorf("unable to decompress log for step %d for build %d: %v", s.GetID(), s.GetBuildID(), err)

		// return the uncompressed log
		return l.ToAPI(), nil
	}

	// return the log
	return l.ToAPI(), nil
}
