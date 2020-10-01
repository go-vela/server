// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

const (
	// CreateWorkerTable represents a query to
	// create the workers table for Vela.
	CreateWorkerTable = `
CREATE TABLE
IF NOT EXISTS
workers (
	id               SERIAL PRIMARY KEY,
	hostname         VARCHAR(250),
	address          VARCHAR(250),
	routes           VARCHAR(1000),
	active           BOOLEAN,
	last_checked_in  INTEGER,
	UNIQUE(hostname)
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
