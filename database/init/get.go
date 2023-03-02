// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package init

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetInit gets a hook by ID from the database.
func (e *engine) GetInit(id int64) (*library.Init, error) {
	e.logger.Tracef("getting hook %d from the database", id)

	// variable to store query results
	i := new(database.Init)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableInit).
		Where("id = ?", id).
		Take(i).
		Error
	if err != nil {
		return nil, err
	}

	// return the hook
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Init.ToLibrary
	return i.ToLibrary(), nil
}
