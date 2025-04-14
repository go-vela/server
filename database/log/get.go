// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetLog gets a log by ID from the database.
func (e *Engine) GetLog(ctx context.Context, id int64) (*api.Log, error) {
	e.logger.Tracef("getting log %d", id)

	// variable to store query results
	l := new(types.Log)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableLog).
		Where("id = ?", id).
		Take(l).
		Error
	if err != nil {
		return nil, err
	}

	// decompress log data
	err = l.Decompress()
	if err != nil {
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allows us to fetch uncompressed logs
		e.logger.Errorf("unable to decompress log %d: %v", id, err)

		// return the uncompressed log
		return l.ToAPI(), nil
	}

	// return the log
	return l.ToAPI(), nil
}
