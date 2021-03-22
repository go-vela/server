// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"

	"github.com/go-vela/server/database"

	"github.com/sirupsen/logrus"
)

// ClientOpt represents a configuration option to initialize the secret client.
type ClientOpt func(*client) error

// WithDatabase sets the Vela database service in the secret client.
func WithDatabase(d database.Service) ClientOpt {
	logrus.Trace("configuring database service in native secret client")

	return func(c *client) error {
		// check if the Vela database service provided is empty
		if d == nil {
			return fmt.Errorf("no Vela database service provided")
		}

		// set the Vela database service in the secret client
		c.Database = d

		return nil
	}
}
