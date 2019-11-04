// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// GetService gets a service by number and build ID from the database.
func (c *client) GetService(number int, b *library.Build) (*library.Service, error) {
	logrus.Tracef("Getting service %d for build %d from the database", number, b.GetNumber())

	// variable to store query results
	s := new(database.Service)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableService).
		Raw(c.DML.ServiceService.Select["build"], b.ID, number).
		Scan(s).Error

	return s.ToLibrary(), err
}

// GetServiceList gets a list of all Services from the database.
func (c *client) GetServiceList() ([]*library.Service, error) {
	logrus.Trace("Listing Services from the database")

	// variable to store query results
	s := new([]database.Service)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableService).
		Raw(c.DML.ServiceService.List["all"]).
		Scan(s).Error

	// variable we want to return
	services := []*library.Service{}
	// iterate through all query results
	for _, service := range *s {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := service

		// convert query result to library type
		services = append(services, tmp.ToLibrary())
	}

	return services, err
}

// GetBuildServiceList gets a list of all services by build ID from the database.
func (c *client) GetBuildServiceList(b *library.Build, page, perPage int) ([]*library.Service, error) {
	logrus.Tracef("Listing services for build %d from the database", b.GetNumber())

	// variable to store query results
	s := new([]database.Service)
	// calculate offset for pagination through results
	offset := (perPage * (page - 1))

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableService).
		Raw(c.DML.ServiceService.List["build"], b.ID, perPage, offset).
		Scan(s).Error

	// variable we want to return
	services := []*library.Service{}
	// iterate through all query results
	for _, service := range *s {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := service

		// convert query result to library type
		services = append(services, tmp.ToLibrary())
	}

	return services, err
}

// GetBuildServiceCount gets a count of all services by build ID from the database.
func (c *client) GetBuildServiceCount(b *library.Build) (int64, error) {
	logrus.Tracef("Counting build services for build %d in the database", b.GetNumber())

	// variable to store query results
	var r []int64

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableService).
		Raw(c.DML.ServiceService.Select["count"], b.ID).
		Pluck("count", &r).Error

	return r[0], err
}

// CreateService creates a new service in the database.
func (c *client) CreateService(s *library.Service) error {
	logrus.Tracef("Creating service %s in the database", s.GetName())

	// cast to database type
	service := database.ServiceFromLibrary(s)

	// validate the necessary fields are populated
	err := service.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Database.
		Table(constants.TableService).
		Create(service).Error
}

// UpdateService updates a service in the database.
func (c *client) UpdateService(s *library.Service) error {
	logrus.Tracef("Updating service %s in the database", s.GetName())

	// cast to database type
	service := database.ServiceFromLibrary(s)

	// validate the necessary fields are populated
	err := service.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Database.
		Table(constants.TableService).
		Where("id = ?", s.ID).
		Update(service).Error
}

// DeleteService deletes a service by unique ID from the database.
func (c *client) DeleteService(id int64) error {
	logrus.Tracef("Deleting service %d from the database", id)

	// send query to the database
	return c.Database.
		Table(constants.TableService).
		Raw(c.DML.ServiceService.Delete, id).Error
}
