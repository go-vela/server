// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package hook

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// ListHooks gets a list of all hooks from the database.
func (e *engine) ListHooks() ([]*library.Hook, error) {
	e.logger.Trace("listing all hooks from the database")

	// variables to store query results and return value
	count := int64(0)
	h := new([]database.Hook)
	hooks := []*library.Hook{}

	// count the results
	count, err := e.CountHooks()
	if err != nil {
		return nil, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return hooks, nil
	}

	// send query to the database and store result in variable
	err = e.client.
		Table(constants.TableHook).
		Find(&h).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, hook := range *h {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := hook

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Hook.ToLibrary
		hooks = append(hooks, tmp.ToLibrary())
	}

	return hooks, nil
}
