// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"

	"github.com/go-vela/server/database"
)

type client struct {
	Native database.Service
}

// New returns a Secret implementation that integrates with a Native secrets engine.
func New(d database.Service) (*client, error) {
	// immediately return if a nil database Service is provided
	if d == nil {
		return nil, fmt.Errorf("empty Database client passed to native secret engine")
	}

	// create the client object
	client := &client{
		Native: d,
	}

	return client, nil
}
