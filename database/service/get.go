// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetService gets a service by ID from the database.
func (e *engine) GetService(ctx context.Context, id int64) (*library.Service, error) {
	e.logger.Tracef("getting service %d", id)

	// variable to store query results
	s := new(database.Service)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableService).
		Where("id = ?", id).
		Take(s).
		Error
	if err != nil {
		return nil, err
	}

	// return the service
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Service.ToLibrary
	return s.ToLibrary(), nil
}
