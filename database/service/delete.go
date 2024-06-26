// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// DeleteService deletes an existing service from the database.
func (e *engine) DeleteService(ctx context.Context, s *library.Service) error {
	e.logger.WithFields(logrus.Fields{
		"service": s.GetNumber(),
	}).Tracef("deleting service %s", s.GetName())

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
