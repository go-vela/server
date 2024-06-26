// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// ListServices gets a list of all services from the database.
func (e *engine) ListServices(ctx context.Context) ([]*library.Service, error) {
	e.logger.Trace("listing all services")

	// variables to store query results and return value
	count := int64(0)
	w := new([]database.Service)
	services := []*library.Service{}

	// count the results
	count, err := e.CountServices(ctx)
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
