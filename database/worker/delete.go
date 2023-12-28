// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// DeleteWorker deletes an existing worker from the database.
func (e *engine) DeleteWorker(ctx context.Context, w *library.Worker) error {
	e.logger.WithFields(logrus.Fields{
		"worker": w.GetHostname(),
	}).Tracef("deleting worker %s from the database", w.GetHostname())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#WorkerFromLibrary
	worker := database.WorkerFromLibrary(w)

	// send query to the database
	return e.client.
		Table(constants.TableWorker).
		Delete(worker).
		Error
}
