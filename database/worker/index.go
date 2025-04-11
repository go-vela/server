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
func (e *Engine) CreateWorkerIndexes(ctx context.Context) error {
	e.logger.Tracef("creating indexes for workers table")

	// create the hostname and address columns index for the workers table
	return e.client.
		WithContext(ctx).
		Exec(CreateHostnameAddressIndex).Error
}
