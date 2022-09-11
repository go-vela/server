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

// GetWorkerForHostname gets a worker by hostname from the database.
func (e *engine) GetWorkerForHostname(hostname string) (*library.Worker, error) {
	e.logger.WithFields(logrus.Fields{
		"worker": hostname,
	}).Tracef("getting worker %s from the database", hostname)

	// variable to store query results
	w := new(database.Worker)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableWorker).
		Where("hostname = ?", hostname).
		Limit(1).
		Take(w).
		Error
	if err != nil {
		return nil, err
	}

	// return the worker
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Worker.ToLibrary
	return w.ToLibrary(), nil
}
