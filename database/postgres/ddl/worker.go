// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
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
	id               SERIAL PRIMARY KEY,
	hostname         VARCHAR(250),
	address          VARCHAR(250),
	routes           VARCHAR(1000),
	active           BOOLEAN,
	last_checked_in  INTEGER,
	build_limit      INTEGER,
	status      VARCHAR(250),
	builds      VARCHAR(250),
	UNIQUE(hostname)
);
`

	// CreateWorkerHostnameAddressIndex represents a query to create an
	// index on the workers table for the hostname and address columns.
	CreateWorkerHostnameAddressIndex = `
CREATE INDEX
IF NOT EXISTS
workers_hostname_address
ON workers (hostname, address);
`
)
