// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package service

import (
	"github.com/go-vela/types/constants"
)

// CountServices gets the count of all services from the database.
func (e *engine) CountServices() (int64, error) {
	e.logger.Tracef("getting count of all services from the database")

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableService).
		Count(&s).
		Error

	return s, err
}