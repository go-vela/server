// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// DeleteLog deletes an existing log from the database.
func (e *Engine) DeleteLog(ctx context.Context, l *api.Log) error {
	// check what the log entry is for
	switch {
	case l.GetServiceID() > 0:
		e.logger.Tracef("deleting log for service %d for build %d", l.GetServiceID(), l.GetBuildID())
	case l.GetStepID() > 0:
		e.logger.Tracef("deleting log for step %d for build %d", l.GetStepID(), l.GetBuildID())
	}

	// cast the API type to database type
	log := types.LogFromAPI(l)

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableLog).
		Delete(log).
		Error
}
