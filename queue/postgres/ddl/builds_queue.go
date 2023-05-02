// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package ddl

const (
	// CreateBuildsQueueTable represents a query to
	// create the builds queue table for Vela.
	CreateBuildsQueueTable = `
CREATE TABLE
IF NOT EXISTS
builds_queue (
	id             	SERIAL PRIMARY KEY,
	channel        	VARCHAR(250),
	payload         BYTEA
);
`

	// CreateBuildsQueueChannelIndex represents a query to create an
	// index on the builds_queue table for the channel column.
	CreateBuildsQueueChannelIndex = `
CREATE INDEX
IF NOT EXISTS
builds_queue_channel
ON builds_queue (channel);
`
)
