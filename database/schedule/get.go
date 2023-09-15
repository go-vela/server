// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	"context"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetSchedule gets a schedule by ID from the database.
func (e *engine) GetSchedule(ctx context.Context, id int64) (*library.Schedule, error) {
	e.logger.Tracef("getting schedule %d from the database", id)

	// variable to store query results
	s := new(database.Schedule)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableSchedule).
		Where("id = ?", id).
		Take(s).
		Error
	if err != nil {
		return nil, err
	}

	return s.ToLibrary(), nil
}
