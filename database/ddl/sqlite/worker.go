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
Workers (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	hostname		TEXT,
	address			TEXT,
	routes          TEXT,
	active          TEXT,
	last_checked_in	INTEGER
);
`

	// CreateWorkersHostnameAddressIndex represents a query to create an
	// index on the repos table for the hostname and address columns.
	CreateWorkersHostnameAddressIndex = `
CREATE INDEX
IF NOT EXISTS
secrets_type_org_repo
ON workers (type, org, repo);
`
)

// createSecretService is a helper function to return
// a service for interacting with the secrets table.
func createWorkerService() *Service {
	return &Service{
		Create:  CreateWorkerTable,
		Indexes: []string{CreateWorkersHostnameAddressIndex},
	}
}
