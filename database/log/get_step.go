// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code with get_service.go
package log

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetLogForStep gets a log by step ID from the database.
func (e *engine) GetLogForStep(ctx context.Context, s *api.Step) (*library.Log, error) {
	e.logger.Tracef("getting log for step %d for build %d", s.GetID(), s.GetBuildID())

	// variable to store query results
	l := new(database.Log)

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
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Log.Decompress
	err = l.Decompress()
	if err != nil {
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allows us to fetch uncompressed logs
		e.logger.Errorf("unable to decompress log for step %d for build %d: %v", s.GetID(), s.GetBuildID(), err)

		// return the uncompressed log
		return l.ToLibrary(), nil
	}

	// return the log
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Log.ToLibrary
	return l.ToLibrary(), nil
}
