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

// GetWorker gets a worker by hostname from the database.
func (c *client) GetWorker(hostname string) (*library.Worker, error) {
	logrus.Tracef("getting worker %s from the database", hostname)

	// variable to store query results
	w := new(database.Worker)

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableWorker).
		Raw(dml.SelectWorker, hostname).
		Scan(w).Error

	return w.ToLibrary(), err
}

// GetWorker gets a worker by address from the database.
func (c *client) GetWorkerByAddress(address string) (*library.Worker, error) {
	logrus.Tracef("getting worker %s from the database", address)

	// variable to store query results
	w := new(database.Worker)

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableWorker).
		Raw(dml.SelectWorkerByAddress, address).
		Scan(w).Error

	return w.ToLibrary(), err
}

// CreateWorker creates a new worker in the database.
func (c *client) CreateWorker(w *library.Worker) error {
	logrus.Tracef("creating worker %s in the database", w.GetHostname())

	// cast to database type
	worker := database.WorkerFromLibrary(w)

	// validate the necessary fields are populated
	err := worker.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Sqlite.
		Table(constants.TableWorker).
		Create(worker).Error
}

// UpdateWorker updates a worker in the database.
func (c *client) UpdateWorker(w *library.Worker) error {
	logrus.Tracef("updating worker %s in the database", w.GetHostname())

	// cast to database type
	worker := database.WorkerFromLibrary(w)

	// validate the necessary fields are populated
	err := worker.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Sqlite.
		Table(constants.TableWorker).
		Save(worker).Error
}

// DeleteWorker deletes a worker by unique ID from the database.
func (c *client) DeleteWorker(id int64) error {
	logrus.Tracef("deleting worker %d in the database", id)

	// send query to the database
	return c.Sqlite.
		Table(constants.TableWorker).
		Exec(dml.DeleteWorker, id).Error
}
