// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// UpdateWorker updates an existing worker in the database.
func (e *engine) UpdateWorker(ctx context.Context, w *api.Worker) (*api.Worker, error) {
	e.logger.WithFields(logrus.Fields{
		"worker": w.GetHostname(),
	}).Tracef("updating worker %s in the database", w.GetHostname())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#WorkerFromLibrary
	worker := WorkerFromAPI(w)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Worker.Validate
	err := worker.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	result := e.client.Table(constants.TableWorker).Save(worker)

	return worker.ToAPI(w.GetRunningBuilds()), result.Error
}
