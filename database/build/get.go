// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetBuild gets a build by ID from the database.
func (e *engine) GetBuild(ctx context.Context, id int64) (*library.Build, error) {
	e.logger.Tracef("getting build %d from the database", id)

	// variable to store query results
	b := new(database.Build)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableBuild).
		Where("id = ?", id).
		Take(b).
		Error
	if err != nil {
		return nil, err
	}

	return b.ToLibrary(), nil
}
