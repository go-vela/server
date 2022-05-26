// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"errors"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database/postgres/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"gorm.io/gorm"
)

// GetWorker gets a worker by hostname from the database.
//
// nolint: dupl // ignore similar code with build
func (c *client) GetWorker(hostname string) (*library.Worker, error) {
	c.Logger.WithFields(logrus.Fields{
		"worker": hostname,
	}).Tracef("getting worker %s from the database", hostname)

	// variable to store query results
	w := new(database.Worker)

	// send query to the database and store result in variable
	result := c.Postgres.
		Table(constants.TableWorker).
		Raw(dml.SelectWorker, hostname).
		Scan(w)

	// check if the query returned a record not found error or no rows were returned
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return w.ToLibrary(), result.Error
}

// GetWorker gets a worker by address from the database.
func (c *client) GetWorkerByAddress(address string) (*library.Worker, error) {
	c.Logger.Tracef("getting worker by address %s from the database", address)

	// variable to store query results
	w := new(database.Worker)

	// send query to the database and store result in variable
	result := c.Postgres.
		Table(constants.TableWorker).
		Raw(dml.SelectWorkerByAddress, address).
		Scan(w)

	// check if the query returned a record not found error or no rows were returned
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return w.ToLibrary(), result.Error
}

// CreateWorker creates a new worker in the database.
func (c *client) CreateWorker(w *library.Worker) error {
	c.Logger.WithFields(logrus.Fields{
		"worker": w.GetHostname(),
	}).Tracef("creating worker %s in the database", w.GetHostname())

	// cast to database type
	worker := database.WorkerFromLibrary(w)

	// validate the necessary fields are populated
	err := worker.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Postgres.
		Table(constants.TableWorker).
		Create(worker).Error
}

// UpdateWorker updates a worker in the database.
func (c *client) UpdateWorker(w *library.Worker) error {
	c.Logger.WithFields(logrus.Fields{
		"worker": w.GetHostname(),
	}).Tracef("updating worker %s in the database", w.GetHostname())

	// cast to database type
	worker := database.WorkerFromLibrary(w)

	// validate the necessary fields are populated
	err := worker.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Postgres.
		Table(constants.TableWorker).
		Save(worker).Error
}

// DeleteWorker deletes a worker by unique ID from the database.
func (c *client) DeleteWorker(id int64) error {
	c.Logger.Tracef("deleting worker %d in the database", id)

	// send query to the database
	return c.Postgres.
		Table(constants.TableWorker).
		Exec(dml.DeleteWorker, id).Error
}
