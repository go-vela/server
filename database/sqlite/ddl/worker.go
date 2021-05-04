// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package ddl

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
	last_checked_in	INTEGER,
	build_limit     INTEGER,
	UNIQUE(hostname)
);
`

	// CreateWorkersHostnameAddressIndex represents a query to create an
	// index on the workers table for the hostname and address columns.
	CreateWorkerHostnameAddressIndex = `
CREATE INDEX
IF NOT EXISTS
workers_hostname_address
ON workers (hostname, address);
`
)
