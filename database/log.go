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

// GetBuildLogs gets a collection of logs for a build by unique ID from the database.
func (c *client) GetBuildLogs(id int64) ([]*library.Log, error) {
	logrus.Tracef("Listing logs for build %d from the database", id)

	// variable to store query results
	l := new([]database.Log)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableLog).
		Raw(c.DML.LogService.List["build"], id).
		Scan(l).Error

	// variable we want to return
	logs := []*library.Log{}
	// iterate through all query results
	for _, log := range *l {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := log

		// convert query result to library type
		logs = append(logs, tmp.ToLibrary())
	}

	return logs, err
}

// GetStepLog gets a log by unique ID from the database.
func (c *client) GetStepLog(id int64) (*library.Log, error) {
	logrus.Tracef("Getting log for step %d from the database", id)

	// variable to store query results
	l := new(database.Log)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableLog).
		Raw(c.DML.LogService.Select["step"], id).
		Scan(l).Error

	return l.ToLibrary(), err
}

// GetServiceLog gets a log by unique ID from the database.
func (c *client) GetServiceLog(id int64) (*library.Log, error) {
	logrus.Tracef("Getting log for service %d from the database", id)

	// variable to store query results
	l := new(database.Log)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableLog).
		Raw(c.DML.LogService.Select["service"], id).
		Scan(l).Error

	return l.ToLibrary(), err
}

// CreateLog creates a new log in the database.
func (c *client) CreateLog(l *library.Log) error {
	logrus.Tracef("Creating log for step %d in the database", l.GetStepID())

	// cast to database type
	log := database.LogFromLibrary(l)

	// validate the necessary fields are populated
	err := log.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Database.
		Table(constants.TableLog).
		Create(log).Error
}

// UpdateLog updates a log in the database.
func (c *client) UpdateLog(l *library.Log) error {
	logrus.Tracef("Updating log for step %d in the database", l.GetStepID())

	// cast to database type
	log := database.LogFromLibrary(l)

	// validate the necessary fields are populated
	err := log.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Database.
		Table(constants.TableLog).
		Where("id = ?", l.GetID()).
		Update(log).Error
}

// DeleteLog deletes a log by unique ID from the database.
func (c *client) DeleteLog(id int64) error {
	logrus.Tracef("Deleting log %d from the database", id)

	// send query to the database
	return c.Database.
		Table(constants.TableLog).
		Raw(c.DML.LogService.Delete, id).Error
}
