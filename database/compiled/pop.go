// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package compiled

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"gorm.io/gorm/clause"
)

// GetPipeline gets a pipeline by ID from the database.
func (e *engine) PopCompiled(id int64) (*library.Compiled, error) {
	e.logger.Tracef("getting pipeline %d from the database", id)

	// variable to store query results
	c := new(database.Compiled)

	switch e.config.Driver {
	case constants.DriverPostgres:
		// send query to the database and store result in variable
		err := e.client.
			Table(constants.TableCompiled).
			Clauses(clause.Returning{}).
			Where("build_id = ?", id).
			Delete(c).
			Error

		if err != nil {
			return nil, err
		}

	case constants.DriverSqlite:
		// send query to the database and store result in variable
		err := e.client.
			Table(constants.TableCompiled).
			Where("id = ?", id).
			Take(c).
			Error
		if err != nil {
			return nil, err
		}

		err = e.client.
			Table(constants.TableCompiled).
			Delete(c).
			Error
		if err != nil {
			return nil, err
		}
	}

	// decompress data for the pipeline
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Pipeline.Decompress
	err := c.Decompress()
	if err != nil {
		return nil, err
	}

	// return the decompressed pipeline
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Pipeline.ToLibrary
	return c.ToLibrary(), nil
}
