// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package log

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// ListLogsForBuild gets a list of logs by build ID from the database.
func (e *engine) ListLogsForBuild(b *library.Build, page, perPage int) ([]*library.Log, int64, error) {
	e.logger.Tracef("listing logs for build %d from the database", b.GetID())

	// variables to store query results and return value
	count := int64(0)
	l := new([]database.Log)
	logs := []*library.Log{}

	// count the results
	count, err := e.CountLogsForBuild(b)
	if err != nil {
		return nil, 0, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return logs, 0, nil
	}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	// send query to the database and store result in variable
	err = e.client.
		Table(constants.TableLog).
		Where("build_id = ?", b.GetID()).
		Order("step_id ASC").
		Limit(perPage).
		Offset(offset).
		Find(&l).
		Error
	if err != nil {
		return nil, count, err
	}

	// iterate through all query results
	for _, log := range *l {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := log

		// decompress log data for the build
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Log.Decompress
		err = tmp.Decompress()
		if err != nil {
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch uncompressed logs
			e.logger.Errorf("unable to decompress logs for build %d: %v", b.GetID(), err)
		}

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Log.ToLibrary
		logs = append(logs, tmp.ToLibrary())
	}

	return logs, count, nil
}
