// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// GetWorker gets a worker by hostname from the database.
func (c *client) GetWorker(hostname string) (*library.Worker, error) {
	logrus.Tracef("Getting worker %s from the database", hostname)

	// variable to store query results
	w := new(database.Worker)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableWorker).
		Raw(c.DML.WorkerService.Select["worker"], hostname).
		Scan(w).Error

	return w.ToLibrary(), err
}

// GetWorker gets a worker by hostname from the database.
func (c *client) GetWorkerByAddress(address string) (*library.Worker, error) {
	logrus.Tracef("Getting worker %s from the database", address)

	// variable to store query results
	w := new(database.Worker)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableWorker).
		Raw(c.DML.WorkerService.Select["address"], address).
		Scan(w).Error

	return w.ToLibrary(), err
}

// GetWorkerList gets a list of all workers from the database.
func (c *client) GetWorkerList() ([]*library.Worker, error) {
	logrus.Trace("Listing workers from the database")

	// variable to store query results
	w := new([]database.Worker)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableWorker).
		Raw(c.DML.WorkerService.List["all"]).
		Scan(w).Error

	// variable we want to return
	workers := []*library.Worker{}
	// iterate through all query results
	for _, worker := range *w {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := worker

		// convert query result to library type
		workers = append(workers, tmp.ToLibrary())
	}

	return workers, err
}

// GetWorkerCount gets a count of all workers from the database.
func (c *client) GetWorkerCount() (int64, error) {
	logrus.Trace("Counting workers in the database")

	// variable to store query results
	var w []int64

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableWorker).
		Raw(c.DML.WorkerService.Select["count"]).
		Pluck("count", &w).Error

	if err != nil {
		return 0, err
	}

	return w[0], nil
}

// CreateWorker creates a new worker in the database.
func (c *client) CreateWorker(w *library.Worker) error {
	logrus.Tracef("Creating worker %s in the database", w.GetHostname())

	// cast to database type
	worker := database.WorkerFromLibrary(w)

	// validate the necessary fields are populated
	err := worker.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Database.
		Table(constants.TableWorker).
		Create(worker).Error
}

// UpdateWorker updates a worker in the database.
func (c *client) UpdateWorker(w *library.Worker) error {
	logrus.Tracef("Updating worker %s in the database", w.GetHostname())

	// cast to database type
	worker := database.WorkerFromLibrary(w)

	// validate the necessary fields are populated
	err := worker.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Database.
		Table(constants.TableWorker).
		Where("id = ?", w.GetID()).
		Save(worker).Error
}

// DeleteWorker deletes a worker by unique ID from the database.
func (c *client) DeleteWorker(id int64) error {
	logrus.Tracef("Deleting worker %d in the database", id)

	// send query to the database
	return c.Database.
		Table(constants.TableWorker).
		Exec(c.DML.WorkerService.Delete, id).Error
}
