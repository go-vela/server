// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"fmt"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CreateLog creates a new log in the database.
func (e *Engine) CreateLog(ctx context.Context, l *api.Log) error {
	// check what the log entry is for
	switch {
	case l.GetServiceID() > 0:
		e.logger.Tracef("creating log for service %d for build %d", l.GetServiceID(), l.GetBuildID())
	case l.GetStepID() > 0:
		e.logger.Tracef("creating log for step %d for build %d", l.GetStepID(), l.GetBuildID())
	}

	// cast the API type to database type
	log := types.LogFromAPI(l)

	err := log.Validate()
	if err != nil {
		return err
	}

	// compress log data for the resource
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
		WithContext(ctx).
		Table(constants.TableLog).
		Create(log).
		Error
}
