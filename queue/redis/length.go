// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package redis

import (
	"context"
)

// Length tallies all items present in the configured channels in the queue.
func (c *client) Length(ctx context.Context) (int64, error) {
	c.Logger.Tracef("reading length of all configured channels in queue")

	total := int64(0)

	for _, channel := range c.config.Channels {
		items, err := c.Redis.LLen(ctx, channel).Result()
		if err != nil {
			return 0, err
		}

		total += items
	}

	return total, nil
}
