// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListServices gets a list of all services from the database.
func (e *Engine) ListServices(ctx context.Context) ([]*api.Service, error) {
	e.logger.Trace("listing all services")

	// variables to store query results and return value
	w := new([]types.Service)
	services := []*api.Service{}

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
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

		services = append(services, tmp.ToAPI())
	}

	return services, nil
}
