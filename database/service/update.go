// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// UpdateService updates an existing service in the database.
func (e *engine) UpdateService(ctx context.Context, s *library.Service) (*library.Service, error) {
	e.logger.WithFields(logrus.Fields{
		"service": s.GetNumber(),
	}).Tracef("updating service %s in the database", s.GetName())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#ServiceFromLibrary
	service := database.ServiceFromLibrary(s)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Service.Validate
	err := service.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	result := e.client.Table(constants.TableService).Save(service)

	return service.ToLibrary(), result.Error
}
