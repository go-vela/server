// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package log

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// ListLogs gets a list of all logs from the database.
func (e *engine) ListLogs() ([]*library.Log, error) {
	e.logger.Trace("listing all logs from the database")

	// variables to store query results and return value
	count := int64(0)
	l := new([]database.Log)
	logs := []*library.Log{}

	// count the results
	count, err := e.CountLogs()
	if err != nil {
		return nil, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return logs, nil
	}

	// send query to the database and store result in variable
	err = e.client.
		Table(constants.TableLog).
		Find(&l).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, log := range *l {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := log

		// decompress log data
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Log.Decompress
		err = tmp.Decompress()
		if err != nil {
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch uncompressed logs
			e.logger.Errorf("unable to decompress logs: %v", err)
		}

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Log.ToLibrary
		logs = append(logs, tmp.ToLibrary())
	}

	return logs, nil
}
