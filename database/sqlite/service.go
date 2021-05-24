// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"github.com/go-vela/server/database/sqlite/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// GetService gets a service by number and build ID from the database.
func (c *client) GetService(number int, b *library.Build) (*library.Service, error) {
	logrus.Tracef("getting service %d for build %d from the database", number, b.GetNumber())

	// variable to store query results
	s := new(database.Service)

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableService).
		Raw(dml.SelectBuildService, b.GetID(), number).
		Scan(s).Error

	return s.ToLibrary(), err
}

// CreateService creates a new service in the database.
func (c *client) CreateService(s *library.Service) error {
	logrus.Tracef("creating service %s in the database", s.GetName())

	// cast to database type
	service := database.ServiceFromLibrary(s)

	// validate the necessary fields are populated
	err := service.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Sqlite.
		Table(constants.TableService).
		Create(service).Error
}

// UpdateService updates a service in the database.
func (c *client) UpdateService(s *library.Service) error {
	logrus.Tracef("updating service %s in the database", s.GetName())

	// cast to database type
	service := database.ServiceFromLibrary(s)

	// validate the necessary fields are populated
	err := service.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Sqlite.
		Table(constants.TableService).
		Save(service).Error
}

// DeleteService deletes a service by unique ID from the database.
func (c *client) DeleteService(id int64) error {
	logrus.Tracef("deleting service %d from the database", id)

	// send query to the database
	return c.Sqlite.
		Table(constants.TableService).
		Exec(dml.DeleteService, id).Error
}
