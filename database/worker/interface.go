// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// WorkerInterface represents the Vela interface for worker
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type WorkerInterface interface {
	// Worker Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateWorkerIndexes defines a function that creates the indexes for the workers table.
	CreateWorkerIndexes(context.Context) error
	// CreateWorkerTable defines a function that creates the workers table.
	CreateWorkerTable(context.Context, string) error

	// Worker Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CountWorkers defines a function that gets the count of all workers.
	CountWorkers(context.Context) (int64, error)
	// CreateWorker defines a function that creates a new worker.
	CreateWorker(context.Context, *api.Worker) (*api.Worker, error)
	// DeleteWorker defines a function that deletes an existing worker.
	DeleteWorker(context.Context, *api.Worker) error
	// GetWorker defines a function that gets a worker by ID.
	GetWorker(context.Context, int64) (*api.Worker, error)
	// GetWorkerForHostname defines a function that gets a worker by hostname.
	GetWorkerForHostname(context.Context, string) (*api.Worker, error)
	// ListWorkers defines a function that gets a list of all workers.
	ListWorkers(context.Context, string, int64, int64) ([]*api.Worker, error)
	// UpdateWorker defines a function that updates an existing worker.
	UpdateWorker(context.Context, *api.Worker) (*api.Worker, error)
}
