// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetWorkerForHostname gets a worker by hostname from the database.
func (e *Engine) GetWorkerForHostname(ctx context.Context, hostname string) (*api.Worker, error) {
	e.logger.WithFields(logrus.Fields{
		"worker": hostname,
	}).Tracef("getting worker %s", hostname)

	// variable to store query results
	w := new(types.Worker)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableWorker).
		Where("hostname = ?", hostname).
		Take(w).
		Error
	if err != nil {
		return nil, err
	}

	// return the worker
	return w.ToAPI(convertToBuilds(w.RunningBuildIDs)), nil
}
