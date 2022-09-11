// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

// nolint: dupl // ignore similar code in update.go
package worker

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CreateWorker creates a new worker in the database.
func (e *engine) CreateWorker(w *library.Worker) error {
	e.logger.WithFields(logrus.Fields{
		"worker": w.GetHostname(),
	}).Tracef("creating worker %s in the database", w.GetHostname())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#WorkerFromLibrary
	worker := database.WorkerFromLibrary(w)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Worker.Validate
	err := worker.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return e.client.
		Table(constants.TableWorker).
		Create(worker).
		Error
}
