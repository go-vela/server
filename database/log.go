// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io/ioutil"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// set the compression level for the data stored in the logs.
const logCompressionLevel = 3

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

		// create new buffer from compressed log data
		b := bytes.NewBuffer(tmp.Data)

		// create new reader for reading compressed log data
		r, err := zlib.NewReader(b)
		if err != nil {
			return logs, fmt.Errorf("unable to create new reader: %v", err)
		}

		// defer closing reader
		defer r.Close()

		// capture decompressed log data from compressed log data
		data, err := ioutil.ReadAll(r)
		if err != nil {
			return logs, fmt.Errorf("unable to read log data: %v", err)
		}

		// overwrite database log data with decompressed log data
		tmp.Data = data

		// convert query result to library type
		logs = append(logs, tmp.ToLibrary())
	}

	return logs, err
}

// GetStepLog gets a log by unique ID from the database.
//
// nolint: dupl // ignore similar code with service
func (c *client) GetStepLog(id int64) (*library.Log, error) {
	logrus.Tracef("Getting log for step %d from the database", id)

	// variable to store query results
	l := new(database.Log)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableLog).
		Raw(c.DML.LogService.Select["step"], id).
		Scan(l).Error
	if err != nil {
		return l.ToLibrary(), err
	}

	// create new buffer from compressed log data
	b := bytes.NewBuffer(l.Data)

	// create new reader for reading compressed log data
	r, err := zlib.NewReader(b)
	if err != nil {
		return l.ToLibrary(), fmt.Errorf("unable to create new reader: %v", err)
	}

	// defer closing reader
	defer r.Close()

	// capture decompressed log data from compressed log data
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return l.ToLibrary(), fmt.Errorf("unable to read log data: %v", err)
	}

	// overwrite database log data with decompressed log data
	l.Data = data

	return l.ToLibrary(), err
}

// GetServiceLog gets a log by unique ID from the database.
//
// nolint: dupl // ignore similar code with step
func (c *client) GetServiceLog(id int64) (*library.Log, error) {
	logrus.Tracef("Getting log for service %d from the database", id)

	// variable to store query results
	l := new(database.Log)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableLog).
		Raw(c.DML.LogService.Select["service"], id).
		Scan(l).Error
	if err != nil {
		return l.ToLibrary(), err
	}

	// create new buffer from compressed log data
	b := bytes.NewBuffer(l.Data)

	// create new reader for reading compressed log data
	r, err := zlib.NewReader(b)
	if err != nil {
		return l.ToLibrary(), fmt.Errorf("unable to create new reader: %v", err)
	}

	// defer closing reader
	defer r.Close()

	// capture decompressed log data from compressed log data
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return l.ToLibrary(), fmt.Errorf("unable to read log data: %v", err)
	}

	// overwrite database log data with decompressed log data
	l.Data = data

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

	// create new buffer from storing compressed log data
	b := new(bytes.Buffer)

	// create new writer for writing compressed log data
	w, err := zlib.NewWriterLevel(b, logCompressionLevel)
	if err != nil {
		return err
	}

	// write compressed log data to buffer
	_, err = w.Write(log.Data)
	if err != nil {
		return err
	}

	// close the writer
	//
	// compressed bytes are not flushed until the
	// writer is closed or explicitly flushed
	w.Close()

	// overwrite database log data with compressed log data
	log.Data = b.Bytes()

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

	// create new buffer from storing compressed log data
	b := new(bytes.Buffer)

	// create new writer for writing compressed log data
	w, err := zlib.NewWriterLevel(b, 3)
	if err != nil {
		return err
	}

	// write compressed log data to buffer
	_, err = w.Write(log.Data)
	if err != nil {
		return err
	}

	// closing the writer
	//
	// compressed bytes are not flushed until the
	// writer is closed or explicitly flushed
	w.Close()

	// overwrite database log data with compressed log data
	log.Data = b.Bytes()

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
		Exec(c.DML.LogService.Delete, id).Error
}
