// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
)

// RouteLength returns count of all items present in the given route.
func (c *Client) RouteLength(ctx context.Context, route string) (int64, error) {
	c.Logger.Tracef("reading length of all configured routes in queue")

	items, err := c.Redis.LLen(ctx, route).Result()
	if err != nil {
		return 0, err
	}

	return items, nil
}
