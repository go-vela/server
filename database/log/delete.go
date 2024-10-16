// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// DeleteLog deletes an existing log from the database.
func (e *engine) DeleteLog(ctx context.Context, l *library.Log) error {
	// check what the log entry is for
	switch {
	case l.GetServiceID() > 0:
		e.logger.Tracef("deleting log for service %d for build %d", l.GetServiceID(), l.GetBuildID())
	case l.GetStepID() > 0:
		e.logger.Tracef("deleting log for step %d for build %d", l.GetStepID(), l.GetBuildID())
	}

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#LogFromLibrary
	log := database.LogFromLibrary(l)

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableLog).
		Delete(log).
		Error
}
