// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code with get_step.go
package log

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetLogForService gets a log by service ID from the database.
func (e *engine) GetLogForService(ctx context.Context, s *library.Service) (*library.Log, error) {
	e.logger.Tracef("getting log for service %d for build %d from the database", s.GetID(), s.GetBuildID())

	// variable to store query results
	l := new(database.Log)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableLog).
		Where("service_id = ?", s.GetID()).
		Take(l).
		Error
	if err != nil {
		return nil, err
	}

	// decompress log data for the service
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Log.Decompress
	err = l.Decompress()
	if err != nil {
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allows us to fetch uncompressed logs
		e.logger.Errorf("unable to decompress log for service %d for build %d: %v", s.GetID(), s.GetBuildID(), err)

		// return the uncompressed log
		return l.ToLibrary(), nil
	}

	// return the log
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Log.ToLibrary
	return l.ToLibrary(), nil
}
