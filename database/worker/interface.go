// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"github.com/go-vela/types/library"
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
	CreateWorkerIndexes() error
	// CreateWorkerTable defines a function that creates the workers table.
	CreateWorkerTable(string) error

	// Worker Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CountWorkers defines a function that gets the count of all workers.
	CountWorkers() (int64, error)
	// CreateWorker defines a function that creates a new worker.
	CreateWorker(*library.Worker) error
	// DeleteWorker defines a function that deletes an existing worker.
	DeleteWorker(*library.Worker) error
	// GetWorker defines a function that gets a worker by ID.
	GetWorker(int64) (*library.Worker, error)
	// GetWorkerForHostname defines a function that gets a worker by hostname.
	GetWorkerForHostname(string) (*library.Worker, error)
	// ListWorkers defines a function that gets a list of all workers.
	ListWorkers() ([]*library.Worker, error)
	// UpdateWorker defines a function that updates an existing worker.
	UpdateWorker(*library.Worker) error
}
