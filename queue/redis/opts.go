// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package redis

import (
	"fmt"
	"time"
)

// ClientOpt represents a configuration option to initialize the queue client for Redis.
type ClientOpt func(*client) error

// WithAddress sets the address in the queue client for Redis.
func WithAddress(address string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring address in redis queue client")

		// check if the address provided is empty
		if len(address) == 0 {
			return fmt.Errorf("no Redis queue address provided")
		}

		// set the queue address in the redis client
		c.config.Address = address

		return nil
	}
}

// WithChannels sets the channels in the queue client for Redis.
func WithChannels(channels ...string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring channels in redis queue client")

		// check if the channels provided are empty
		if len(channels) == 0 {
			return fmt.Errorf("no Redis queue channels provided")
		}

		// set the queue channels in the redis client
		c.config.Channels = channels

		return nil
	}
}

// WithCluster sets the clustering mode in the queue client for Redis.
func WithCluster(cluster bool) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring clustering mode in redis queue client")

		// set the queue clustering mode in the redis client
		c.config.Cluster = cluster

		return nil
	}
}

// WithTimeout sets the timeout in the queue client for Redis.
func WithTimeout(timeout time.Duration) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring timeout in redis queue client")

		// set the queue timeout in the redis client
		c.config.Timeout = timeout

		return nil
	}
}
