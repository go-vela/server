// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package service

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// DeleteService deletes an existing service from the database.
func (e *engine) DeleteService(s *library.Service) error {
	e.logger.WithFields(logrus.Fields{
		"service": s.GetNumber(),
	}).Tracef("deleting service %s from the database", s.GetName())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#ServiceFromLibrary
	service := database.ServiceFromLibrary(s)

	// send query to the database
	return e.client.
		Table(constants.TableService).
		Delete(service).
		Error
}
