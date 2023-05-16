// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package service

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// ListServices gets a list of all services from the database.
func (e *engine) ListServices() ([]*library.Service, error) {
	e.logger.Trace("listing all services from the database")

	// variables to store query results and return value
	count := int64(0)
	w := new([]database.Service)
	services := []*library.Service{}

	// count the results
	count, err := e.CountServices()
	if err != nil {
		return nil, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return services, nil
	}

	// send query to the database and store result in variable
	err = e.client.
		Table(constants.TableService).
		Find(&w).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, service := range *w {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := service

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Service.ToLibrary
		services = append(services, tmp.ToLibrary())
	}

	return services, nil
}
