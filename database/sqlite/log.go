// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"errors"
	"fmt"

	"github.com/go-vela/server/database/sqlite/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"gorm.io/gorm"
)

// GetBuildLogs gets a collection of logs for a build by unique ID from the database.
func (c *client) GetBuildLogs(id int64) ([]*library.Log, error) {
	c.Logger.Tracef("listing logs for build %d from the database", id)

	// variable to store query results
	l := new([]database.Log)

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableLog).
		Raw(dml.ListBuildLogs, id).
		Scan(l).Error
	if err != nil {
		return nil, err
	}

	// variable we want to return
	logs := []*library.Log{}
	// iterate through all query results
	for _, log := range *l {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := log

		// decompress log data for the step
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Log.Decompress
		err = tmp.Decompress()
		if err != nil {
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch uncompressed logs
			c.Logger.Errorf("unable to decompress logs for build %d: %v", id, err)
		}

		// convert query result to library type
		logs = append(logs, tmp.ToLibrary())
	}

	return logs, nil
}

// GetStepLog gets a log by unique ID from the database.
//
// nolint: dupl // ignore similar code with service
func (c *client) GetStepLog(id int64) (*library.Log, error) {
	c.Logger.Tracef("getting log for step %d from the database", id)

	// variable to store query results
	l := new(database.Log)

	// send query to the database and store result in variable
	result := c.Sqlite.
		Table(constants.TableLog).
		Raw(dml.SelectStepLog, id).
		Scan(l)

	// check if the query returned a record not found error or no rows were returned
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// decompress log data for the step
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Log.Decompress
	err := l.Decompress()
	if err != nil {
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allows us to fetch uncompressed logs
		c.Logger.Errorf("unable to decompress logs for step %d: %v", id, err)

		// return the uncompressed log
		return l.ToLibrary(), result.Error
	}

	// return the decompressed log
	return l.ToLibrary(), result.Error
}

// GetServiceLog gets a log by unique ID from the database.
//
// nolint: dupl // ignore similar code with step
func (c *client) GetServiceLog(id int64) (*library.Log, error) {
	c.Logger.Tracef("getting log for service %d from the database", id)

	// variable to store query results
	l := new(database.Log)

	// send query to the database and store result in variable
	result := c.Sqlite.
		Table(constants.TableLog).
		Raw(dml.SelectServiceLog, id).
		Scan(l)

	// check if the query returned a record not found error or no rows were returned
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// decompress log data for the service
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Log.Decompress
	err := l.Decompress()
	if err != nil {
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allowing us to fetch uncompressed logs
		c.Logger.Errorf("unable to decompress logs for service %d: %v", id, err)

		// return the uncompressed log
		return l.ToLibrary(), result.Error
	}

	// return the decompressed log
	return l.ToLibrary(), result.Error
}

// CreateLog creates a new log in the database.
//
// nolint: dupl // ignore false positive of duplicate code
func (c *client) CreateLog(l *library.Log) error {
	// check if the log entry is for a step
	if l.GetStepID() > 0 {
		c.Logger.Tracef("creating log for step %d in the database", l.GetStepID())
	} else {
		c.Logger.Tracef("creating log for service %d in the database", l.GetServiceID())
	}

	// cast to database type
	log := database.LogFromLibrary(l)

	// validate the necessary fields are populated
	err := log.Validate()
	if err != nil {
		return err
	}

	// compress log data for the resource
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Log.Compress
	err = log.Compress(c.config.CompressionLevel)
	if err != nil {
		return fmt.Errorf("unable to compress logs for step %d: %v", l.GetStepID(), err)
	}

	// send query to the database
	return c.Sqlite.
		Table(constants.TableLog).
		Create(log).Error
}

// UpdateLog updates a log in the database.
//
// nolint: dupl // ignore false positive of duplicate code
func (c *client) UpdateLog(l *library.Log) error {
	// check if the log entry is for a step
	if l.GetStepID() > 0 {
		c.Logger.Tracef("updating log for step %d in the database", l.GetStepID())
	} else {
		c.Logger.Tracef("updating log for service %d in the database", l.GetServiceID())
	}

	// cast to database type
	log := database.LogFromLibrary(l)

	// validate the necessary fields are populated
	err := log.Validate()
	if err != nil {
		return err
	}

	// compress log data for the resource
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Log.Compress
	err = log.Compress(c.config.CompressionLevel)
	if err != nil {
		return fmt.Errorf("unable to compress logs for step %d: %v", l.GetStepID(), err)
	}

	// send query to the database
	return c.Sqlite.
		Table(constants.TableLog).
		Save(log).Error
}

// DeleteLog deletes a log by unique ID from the database.
func (c *client) DeleteLog(id int64) error {
	c.Logger.Tracef("deleting log %d from the database", id)

	// send query to the database
	return c.Sqlite.
		Table(constants.TableLog).
		Exec(dml.DeleteLog, id).Error
}
