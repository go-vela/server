// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package redis

import (
	"context"

	"github.com/sirupsen/logrus"
)

// Push inserts an item to the specified channel in the queue.
func (c *client) Push(ctx context.Context, channel string, item []byte) error {
	logrus.Tracef("pushing item to queue %s", channel)

	// build a redis queue command to push an item to queue
	//
	// https://pkg.go.dev/github.com/go-redis/redis?tab=doc#Client.RPush
	pushCmd := c.Queue.RPush(ctx, channel, item)

	// blocking call to push an item to queue and return err
	//
	// https://pkg.go.dev/github.com/go-redis/redis?tab=doc#IntCmd.Err
	err := pushCmd.Err()
	if err != nil {
		return err
	}

	return nil
}
