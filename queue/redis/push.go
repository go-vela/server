// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// Push inserts an item to the specified channel in the queue.
func (c *Client) Push(ctx context.Context, channel string, item int64) error {
	c.Logger.Tracef("pushing item to queue %s", channel)

	// ensure the item to be pushed is valid
	// go-redis RPush does not support nil as of v9.0.2
	//
	// https://github.com/redis/go-redis/pull/1960
	if item == 0 {
		return errors.New("item is nil")
	}

	zMember := redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: strconv.FormatInt(item, 10),
	}

	zAddCmd := c.Redis.ZAdd(ctx, channel, zMember)

	err := zAddCmd.Err()
	if err != nil {
		return err
	}

	return nil
}
