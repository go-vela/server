// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package dml

const (
	// ListQueuedBuilds represents a query to
	// list all queued builds in the database.
	ListQueuedBuilds = `
SELECT *
FROM build_queue;
`

	DeleteQueuedBuild = `
DELETE
FROM build_queue
WHERE build_id = ?;
`
)
