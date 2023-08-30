// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// UpdateWorker updates an existing worker in the database.
func (e *engine) UpdateWorker(w *library.Worker) (*library.Worker, error) {
	e.logger.WithFields(logrus.Fields{
		"worker": w.GetHostname(),
	}).Tracef("updating worker %s in the database", w.GetHostname())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#WorkerFromLibrary
	worker := database.WorkerFromLibrary(w)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Worker.Validate
	err := worker.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	result := e.client.Table(constants.TableWorker).Save(worker)

	return worker.ToLibrary(), result.Error
}
