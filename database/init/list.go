// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package init

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// ListInits gets a list of all inits from the database.
func (e *engine) ListInits() ([]*library.Init, error) {
	e.logger.Trace("listing all inits from the database")

	// variables to store query results and return value
	count := int64(0)
	i := new([]database.Init)
	inits := []*library.Init{}

	// count the results
	count, err := e.CountInits()
	if err != nil {
		return nil, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return inits, nil
	}

	// send query to the database and store result in variable
	err = e.client.
		Table(constants.TableInit).
		Find(&i).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, init := range *i {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := init

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Init.ToLibrary
		inits = append(inits, tmp.ToLibrary())
	}

	return inits, nil
}
