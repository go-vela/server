// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// GetWorkerForHostname gets a worker by hostname from the database.
func (e *engine) GetWorkerForHostname(ctx context.Context, hostname string) (*library.Worker, error) {
	e.logger.WithFields(logrus.Fields{
		"worker": hostname,
	}).Tracef("getting worker %s from the database", hostname)

	// variable to store query results
	w := new(database.Worker)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableWorker).
		Where("hostname = ?", hostname).
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
