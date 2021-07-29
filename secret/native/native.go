// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"github.com/go-vela/server/database"
)

// client represents a struct to hold native secret setup.
type client struct {
	// client to interact with database for secret operations
	Database database.Service
}

// New returns a Secret implementation that integrates with a Native secrets engine.
//
// nolint: revive // ignore returning unexported client
func New(opts ...ClientOpt) (*client, error) {
	// create new native client
	c := new(client)

	// create new fields
	c.Database = *new(database.Service)

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}
