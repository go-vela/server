// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/constants"
	"github.com/go-vela/server/database/types"
)

// GetSchedule gets a schedule by ID from the database.
func (e *engine) GetSchedule(id int64) (*api.Schedule, error) {
	e.logger.Tracef("getting schedule %d from the database", id)

	// variable to store query results
	s := new(types.Schedule)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableSchedule).
		Where("id = ?", id).
		Take(s).
		Error
	if err != nil {
		return nil, err
	}

	return s.ToAPI(nil), nil
}
