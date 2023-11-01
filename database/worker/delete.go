// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// DeleteWorker deletes an existing worker from the database.
func (e *engine) DeleteWorker(ctx context.Context, w *api.Worker) error {
	e.logger.WithFields(logrus.Fields{
		"worker": w.GetHostname(),
	}).Tracef("deleting worker %s from the database", w.GetHostname())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#WorkerFromLibrary
	worker := types.WorkerFromAPI(w)

	// send query to the database
	return e.client.
		Table(constants.TableWorker).
		Delete(worker).
		Error
}
