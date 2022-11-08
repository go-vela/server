// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package ddl

const (
	// CreateBuildQueueTable represents a query to
	// create the build queue table for Vela.
	CreateBuildQueueTable = `
CREATE TABLE
IF NOT EXISTS
build_queue (
    id        INTEGER GENERATED ALWAYS AS IDENTITY,
    flavor    VARCHAR(500)      NOT NULL,
	created   INTEGER,
    status    VARCHAR(500),
    full_name VARCHAR(500),
    number    INTEGER,
    build_id  INTEGER,
    pipeline  BYTEA,
    PRIMARY KEY (id)
);
`
)
