// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package redis

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-vela/types"
	"github.com/redis/go-redis/v9"
)

// Pop grabs an item from the specified channel off the queue.
func (c *client) Pop(ctx context.Context) (*types.Item, error) {
	c.Logger.Tracef("popping item from queue %s", c.config.Channels)

	// build a redis queue command to pop an item from queue
	//
	// https://pkg.go.dev/github.com/go-redis/redis?tab=doc#Client.BLPop
	popCmd := c.Redis.BLPop(ctx, c.config.Timeout, c.config.Channels...)

	// blocking call to pop item from queue
	//
	// https://pkg.go.dev/github.com/go-redis/redis?tab=doc#StringSliceCmd.Result
	result, err := popCmd.Result()
	if err != nil {
		// BLPOP timeout
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, err
	}

	item := new(types.Item)

	// unmarshal result into queue item
	err = json.Unmarshal([]byte(result[1]), item)
	if err != nil {
		return nil, err
	}

	return item, nil
}
