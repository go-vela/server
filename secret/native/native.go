// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"

	"github.com/go-vela/server/database"
)

// client represents a struct to hold native secret setup.
type client struct {
	// client to interact with database for secret operations
	Native database.Service
	// key to use for encrypting and decrypting secret values
	passphrase string
}

// New returns a Secret implementation that integrates with a Native secrets engine.
//
// nolint: golint // ignore returning unexported client
func New(d database.Service, passphrase string) (*client, error) {
	// immediately return if a nil database Service is provided
	if d == nil {
		return nil, fmt.Errorf("empty Database client passed to native secret engine")
	}

	// create the client object
	client := &client{
		Native:     d,
		passphrase: passphrase,
	}

	return client, nil
}
