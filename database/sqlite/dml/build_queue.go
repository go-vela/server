// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package dml

const (
	ListQueuedBuilds = `
SELECT *
FROM build_queue;
`
)
