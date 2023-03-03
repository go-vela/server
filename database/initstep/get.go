// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package initstep

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetInitStep gets an init step by ID from the database.
func (e *engine) GetInitStep(id int64) (*library.InitStep, error) {
	e.logger.Tracef("getting init step %d from the database", id)

	// variable to store query results
	i := new(database.InitStep)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableInitStep).
		Where("id = ?", id).
		Take(i).
		Error
	if err != nil {
		return nil, err
	}

	// return the InitStep
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#InitStep.ToLibrary
	return i.ToLibrary(), nil
}
