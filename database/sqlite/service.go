// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"errors"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database/sqlite/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"gorm.io/gorm"
)

// GetService gets a service by number and build ID from the database.
//
//nolint:dupl // ignore similar code with step.
func (c *client) GetService(number int, b *library.Build) (*library.Service, error) {
	c.Logger.WithFields(logrus.Fields{
		"build":   b.GetNumber(),
		"service": number,
	}).Tracef("getting service %d for build %d from the database", number, b.GetNumber())

	// variable to store query results
	s := new(database.Service)

	// send query to the database and store result in variable
	result := c.Sqlite.
		Table(constants.TableService).
		Raw(dml.SelectBuildService, b.GetID(), number).
		Scan(s)

	// check if the query returned a record not found error or no rows were returned
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return s.ToLibrary(), result.Error
}

// CreateService creates a new service in the database.
//
//nolint:dupl // ignore similar code with step.
func (c *client) CreateService(s *library.Service) (*library.Service, error) {
	c.Logger.WithFields(logrus.Fields{
		"service": s.GetNumber(),
	}).Tracef("creating service %s in the database", s.GetName())

	// cast to database type
	service := database.ServiceFromLibrary(s)

	// validate the necessary fields are populated
	err := service.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	err = c.Sqlite.
		Table(constants.TableService).
		Create(service).Error

	if err != nil {
		return nil, err
	}

	return service.ToLibrary(), nil
}

// UpdateService updates a service in the database.
//
//nolint:dupl // ignore similar code with step.
func (c *client) UpdateService(s *library.Service) (*library.Service, error) {
	c.Logger.WithFields(logrus.Fields{
		"service": s.GetNumber(),
	}).Tracef("updating service %s in the database", s.GetName())

	// cast to database type
	service := database.ServiceFromLibrary(s)

	// validate the necessary fields are populated
	err := service.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	err = c.Sqlite.
		Table(constants.TableService).
		Save(service).Error

	if err != nil {
		return nil, err
	}

	return service.ToLibrary(), nil
}

// DeleteService deletes a service by unique ID from the database.
func (c *client) DeleteService(id int64) error {
	c.Logger.Tracef("deleting service %d from the database", id)

	// send query to the database
	return c.Sqlite.
		Table(constants.TableService).
		Exec(dml.DeleteService, id).Error
}
