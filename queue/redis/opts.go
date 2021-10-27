// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package redis

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// ClientOpt represents a configuration option to initialize the queue client.
type ClientOpt func(*client) error

// WithAddress sets the Redis address in the queue client.
func WithAddress(address string) ClientOpt {
	logrus.Trace("configuring address in redis queue client")

	return func(c *client) error {
		// check if the address provided is empty
		if len(address) == 0 {
			return fmt.Errorf("no Redis queue address provided")
		}

		// set the queue address in the redis client
		c.config.Address = address

		return nil
	}
}

// WithChannels sets the Redis channels in the queue client.
func WithChannels(channels ...string) ClientOpt {
	logrus.Trace("configuring channels in redis queue client")

	return func(c *client) error {
		// check if the channels provided are empty
		if len(channels) == 0 {
			return fmt.Errorf("no Redis queue channels provided")
		}

		// set the queue channels in the redis client
		c.config.Channels = channels

		return nil
	}
}

// WithCluster sets the Redis clustering mode in the queue client.
func WithCluster(cluster bool) ClientOpt {
	logrus.Trace("configuring clustering mode in redis queue client")

	return func(c *client) error {
		// set the queue clustering mode in the redis client
		c.config.Cluster = cluster

		return nil
	}
}

// WithTimeout sets the Redis timeout in the queue client.
func WithTimeout(timeout time.Duration) ClientOpt {
	logrus.Trace("configuring timeout in redis queue client")

	return func(c *client) error {
		// set the queue timeout in the redis client
		c.config.Timeout = timeout

		return nil
	}
}
