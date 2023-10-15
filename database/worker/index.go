// SPDX-License-Identifier: Apache-2.0

package worker

import "context"

const (
	// CreateHostnameAddressIndex represents a query to create an
	// index on the workers table for the hostname and address columns.
	CreateHostnameAddressIndex = `
CREATE INDEX
IF NOT EXISTS
workers_hostname_address
ON workers (hostname, address);
`
)

// CreateWorkerIndexes creates the indexes for the workers table in the database.
func (e *engine) CreateWorkerIndexes(ctx context.Context) error {
	e.logger.Tracef("creating indexes for workers table in the database")

	// create the hostname and address columns index for the workers table
	return e.client.Exec(CreateHostnameAddressIndex).Error
}
