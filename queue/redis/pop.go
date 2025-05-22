// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

// Pop grabs an item from the specified channel off the queue.
func (c *Client) Pop(ctx context.Context, inRoutes []string) (int64, error) {
	// define routes to pop from
	var routes []string

	// if routes were supplied, use those
	if len(inRoutes) > 0 {
		routes = inRoutes
	} else {
		routes = c.GetRoutes()
	}

	c.Logger.Tracef("popping item from queue %s", routes)

	zPopCmd := c.Redis.BZPopMin(ctx, c.config.Timeout, routes...)

	result, err := zPopCmd.Result()
	if err != nil {
		// BLPOP timeout
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}

		return 0, err
	}

	bID, ok := result.Member.(int64)
	if !ok {
		return 0, errors.New("failed to convert item to int64")
	}

	return bID, nil
}
