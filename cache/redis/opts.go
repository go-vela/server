// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"errors"
	"fmt"
)

// ClientOpt represents a configuration option to initialize the queue client for Redis.
type ClientOpt func(*Client) error

// WithAddress sets the address in the queue client for Redis.
func WithAddress(address string) ClientOpt {
	return func(c *Client) error {
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

// WithCluster sets the clustering mode in the queue client for Redis.
func WithCluster(cluster bool) ClientOpt {
	return func(c *Client) error {
		c.Logger.Trace("configuring clustering mode in redis queue client")

		// set the queue clustering mode in the redis client
		c.config.Cluster = cluster

		return nil
	}
}

// WithInstallTokenKey sets the install token key in the cache client for Redis.
func WithInstallTokenKey(key string) ClientOpt {
	return func(c *Client) error {
		c.Logger.Trace("configuring install token key in redis cache client")

		// check if the install token key provided is empty
		if len(key) == 0 {
			return errors.New("no install token key provided")
		}

		// set the install token key in the redis cache client
		c.config.InstallTokenKey = key

		return nil
	}
}
