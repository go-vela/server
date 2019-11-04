// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

const (
	// ListLogs represents a query to
	// list all logs in the database.
	ListLogs = `
SELECT *
FROM logs;
`

	// ListBuildLogs represents a query to list
	// all logs for a build_id in the database.
	ListBuildLogs = `
SELECT *
FROM logs
WHERE build_id = $1
ORDER BY step_id ASC;
`

	// SelectStepLog represents a query to select
	// a log for a step_id in the database.
	SelectStepLog = `
SELECT *
FROM logs
WHERE step_id = $1
LIMIT 1;
`

	// SelectServiceLog represents a query to select
	// a log for a service_id in the database.
	SelectServiceLog = `
SELECT *
FROM logs
WHERE service_id = $1
LIMIT 1;
`

	// DeleteLog represents a query to
	// remove a log from the database.
	DeleteLog = `
DELETE
FROM logs
WHERE id = $1
LIMIT 1;
`
)

// createLogService is a helper function to return
// a service for interacting with the logs table.
func createLogService() *Service {
	return &Service{
		List: map[string]string{
			"all":   ListLogs,
			"build": ListBuildLogs,
		},
		Select: map[string]string{
			"step":    SelectStepLog,
			"service": SelectServiceLog,
		},
		Delete: DeleteLog,
	}
}
