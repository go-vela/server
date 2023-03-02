// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package initstep

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// ListInitSteps gets a list of all inits from the database.
func (e *engine) ListInitSteps() ([]*library.InitStep, error) {
	e.logger.Trace("listing all init steps from the database")

	// variables to store query results and return value
	count := int64(0)
	i := new([]database.InitStep)
	initSteps := []*library.InitStep{}

	// count the results
	count, err := e.CountInitSteps()
	if err != nil {
		return nil, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return initSteps, nil
	}

	// send query to the database and store result in variable
	err = e.client.
		Table(constants.TableInitStep).
		Find(&i).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, initStep := range *i {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := initStep

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#InitStep.ToLibrary
		initSteps = append(initSteps, tmp.ToLibrary())
	}

	return initSteps, nil
}
