// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
)

// Length tallies all items present in the configured routes in the queue.
func (c *Client) Length(ctx context.Context) (int64, error) {
	c.Logger.Tracef("reading length of all configured routes in queue")

	total := int64(0)

	for _, route := range c.GetRoutes() {
		items, err := c.Redis.LLen(ctx, route).Result()
		if err != nil {
			return 0, err
		}

		total += items
	}

	return total, nil
}
