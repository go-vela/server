// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetLog gets a log by ID from the database.
func (e *engine) GetLog(ctx context.Context, id int64) (*library.Log, error) {
	e.logger.Tracef("getting log %d", id)

	// variable to store query results
	l := new(database.Log)

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
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Log.Decompress
	err = l.Decompress()
	if err != nil {
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allows us to fetch uncompressed logs
		e.logger.Errorf("unable to decompress log %d: %v", id, err)

		// return the uncompressed log
		return l.ToLibrary(), nil
	}

	// return the log
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Log.ToLibrary
	return l.ToLibrary(), nil
}
