// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code with update.go
package log

import (
	"context"
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// UpdateLog updates an existing log in the database.
func (e *engine) UpdateLog(ctx context.Context, l *library.Log) error {
	// check what the log entry is for
	switch {
	case l.GetServiceID() > 0:
		e.logger.Tracef("updating log for service %d for build %d", l.GetServiceID(), l.GetBuildID())
	case l.GetStepID() > 0:
		e.logger.Tracef("updating log for step %d for build %d", l.GetStepID(), l.GetBuildID())
	}

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#LogFromLibrary
	log := database.LogFromLibrary(l)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Log.Validate
	err := log.Validate()
	if err != nil {
		return err
	}

	// compress log data for the resource
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Log.Compress
	err = log.Compress(e.config.CompressionLevel)
	if err != nil {
		switch {
		case l.GetServiceID() > 0:
			return fmt.Errorf("unable to compress log for service %d for build %d: %w", l.GetServiceID(), l.GetBuildID(), err)
		case l.GetStepID() > 0:
			return fmt.Errorf("unable to compress log for step %d for build %d: %w", l.GetStepID(), l.GetBuildID(), err)
		}
	}

	// send query to the database
	return e.client.
		Table(constants.TableLog).
		Save(log).
		Error
}
