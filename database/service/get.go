// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package service

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetService gets a service by ID from the database.
func (e *engine) GetService(ctx context.Context, id int64) (*library.Service, error) {
	e.logger.Tracef("getting service %d from the database", id)

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
