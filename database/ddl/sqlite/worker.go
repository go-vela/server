// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

const (
	// CreateWorkerTable represents a query to
	// create the workers table for Vela.
	CreateWorkerTable = `
CREATE TABLE
IF NOT EXISTS
workers (
	id              INTEGER PRIMARY KEY AUTOINCREMENT,
	hostname        TEXT,
	address         TEXT,
	routes          TEXT,
	active          TEXT,
	last_checked_in	INTEGER
);
`

	// CreateWorkersHostnameAddressIndex represents a query to create an
	// index on the workers table for the hostname and address columns.
	CreateWorkersHostnameAddressIndex = `
CREATE INDEX
IF NOT EXISTS
workers_hostname_address
ON workers (hostname, address);
`
)

// createWorkerService is a helper function to return
// a service for interacting with the workers table.
func createWorkerService() *Service {
	return &Service{
		Create:  CreateWorkerTable,
		Indexes: []string{CreateWorkersHostnameAddressIndex},
	}
}
