// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package dml

const (
	// ListWorkers represents a query to
	// list all workers in the database.
	ListWorkers = `
SELECT *
FROM workers;
`

	// SelectWorkersCount represents a query to select the
	// count of workers in the database.
	SelectWorkersCount = `
SELECT count(*) as count
FROM workers;
`

	// SelectWorker represents a query to select a
	// worker by hostname in the database.
	SelectWorker = `
SELECT *
FROM workers
WHERE hostname = ?
LIMIT 1;
`

	// SelectWorkerByAddress represents a query to select a
	// worker by address in the database.
	SelectWorkerByAddress = `
SELECT *
FROM workers
WHERE address = ?
LIMIT 1;
`

	// DeleteWorker represents a query to
	// remove a worker from the database.
	DeleteWorker = `
DELETE
FROM workers
WHERE id = ?;
`
)
