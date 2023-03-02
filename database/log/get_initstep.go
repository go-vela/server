// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

//nolint:dupl // ignore similar code with get_step.go
package log

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetLogForInitStep gets a log by init step ID from the database.
func (e *engine) GetLogForInitStep(i *library.InitStep) (*library.Log, error) {
	e.logger.Tracef("getting log for init step %d for build %d from the database", i.GetID(), i.GetBuildID())

	// variable to store query results
	l := new(database.Log)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableLog).
		Where("initstep_id = ?", i.GetID()).
		Take(l).
		Error
	if err != nil {
		return nil, err
	}

	// decompress log data for the init step
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Log.Decompress
	err = l.Decompress()
	if err != nil {
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allowing us to fetch uncompressed logs
		e.logger.Errorf("unable to decompress log for init step %d for build %d: %v", i.GetID(), i.GetBuildID(), err)

		// return the uncompressed log
		return l.ToLibrary(), nil
	}

	// return the log
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Log.ToLibrary
	return l.ToLibrary(), nil
}
