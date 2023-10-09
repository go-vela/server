// SPDX-License-Identifier: Apache-2.0

package native

import (
	"fmt"

	"github.com/go-vela/server/database"
)

// ClientOpt represents a configuration option to initialize the secret client for Native.
type ClientOpt func(*client) error

// WithDatabase sets the Vela database service in the secret client for Native.
func WithDatabase(d database.Interface) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring database service in native secret client")

		// check if the Vela database service provided is empty
		if d == nil {
			return fmt.Errorf("no Vela database service provided")
		}

		// set the Vela database service in the secret client
		c.Database = d

		return nil
	}
}
