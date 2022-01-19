// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
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

// GetServiceList gets a list of all services from the database.
func (c *client) GetServiceList() ([]*library.Service, error) {
	c.Logger.Trace("listing services from the database")

	// variable to store query results
	s := new([]database.Service)

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableService).
		Raw(dml.ListServices).
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

// GetBuildServiceList gets a list of services by build ID from the database.
//
// nolint: lll // ignore long line length due to parameters
func (c *client) GetBuildServiceList(b *library.Build, page, perPage int) ([]*library.Service, error) {
	c.Logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("listing services for build %d from the database", b.GetNumber())

	// variable to store query results
	s := new([]database.Service)
	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableService).
		Raw(dml.ListBuildServices, b.GetID(), perPage, offset).
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
