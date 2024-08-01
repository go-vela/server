// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
)

// Ping contacts the queue to test its connection.
func (c *client) Ping(ctx context.Context) error {
	// send ping request to client
	err := c.Redis.Ping(ctx).Err()
	if err != nil {
		c.Logger.Debugf("unable to ping Redis queue.")
		return err
	}

	return nil
}
