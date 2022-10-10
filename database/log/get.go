// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package log

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetLog gets a log by ID from the database.
func (e *engine) GetLog(id int64) (*library.Log, error) {
	e.logger.Tracef("getting log %d from the database", id)

	// variable to store query results
	l := new(database.Log)

	// send query to the database and store result in variable
	err := e.client.
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
		// which allowing us to fetch uncompressed logs
		e.logger.Errorf("unable to decompress log %d: %v", id, err)

		// return the uncompressed log
		return l.ToLibrary(), nil
	}

	// return the log
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Log.ToLibrary
	return l.ToLibrary(), nil
}
