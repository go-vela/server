// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package itinerary

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"gorm.io/gorm/clause"
)

// PopBuildItinerary pops a build itinerary by build_id from the database.
func (e *engine) PopBuildItinerary(id int64) (*library.BuildItinerary, error) {
	e.logger.Tracef("popping build itinerary for build %d from the database", id)

	// variable to store query results
	b := new(database.BuildItinerary)

	// at the time of coding, GORM does not implement a version of Sqlite3 that supports RETURNING.
	// so we have to select and delete for the Sqlite driver.
	switch e.config.Driver {
	case constants.DriverPostgres:
		// send query to the database and store result in variable
		err := e.client.
			Table(constants.TableBuildItinerary).
			Clauses(clause.Returning{}).
			Where("build_id = ?", id).
			Delete(b).
			Error

		if err != nil {
			return nil, err
		}

	case constants.DriverSqlite:
		// send query to the database and store result in variable
		err := e.client.
			Table(constants.TableBuildItinerary).
			Where("id = ?", id).
			Take(b).
			Error
		if err != nil {
			return nil, err
		}

		// send query to the database to delete result just got
		err = e.client.
			Table(constants.TableBuildItinerary).
			Delete(b).
			Error
		if err != nil {
			return nil, err
		}
	}

	// decompress data for the build itinerary
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#BuildItinerary.Decompress
	err := b.Decompress()
	if err != nil {
		return nil, err
	}

	// return the decompressed build itinerary
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#BuildItinerary.ToLibrary
	return b.ToLibrary(), nil
}
