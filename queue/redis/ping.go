// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"fmt"
)

// Ping contacts the queue to test its connection.
func (c *client) Ping(ctx context.Context) error {
	// send ping request to client
	err := c.Redis.Ping(ctx).Err()
	if err != nil {
		c.Logger.Debugf("unable to ping Redis queue.")
		return fmt.Errorf("unable to establish connection to Redis queue")
	}

	return nil
}
