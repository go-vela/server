// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package redis

// Publish inserts an item to the specified channel in the queue.
func (c *client) Publish(channel string, item []byte) error {
	return c.Queue.RPush(channel, item).Err()
}
