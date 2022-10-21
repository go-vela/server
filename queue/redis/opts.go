// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package redis

import (
	"encoding/base64"
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

// WithPrivateKey sets the private key in the queue client for Redis.
func WithPrivateKey(privateKeyEncoded string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring private key in redis queue client")

		privateKeyDecoded, err := base64.StdEncoding.DecodeString(privateKeyEncoded)
		if err != nil {
			return err
		}

		c.config.SigningPrivateKey = new([64]byte)
		copy(c.config.SigningPrivateKey[:], privateKeyDecoded)

		return nil
	}
}

// WithPublicKey sets the public key in the queue client for Redis.
func WithPublicKey(publicKeyEncoded string) ClientOpt {
	return func(c *client) error {
		c.Logger.Tracef("configuring public key in redis queue client")

		publicKeyDecoded, err := base64.StdEncoding.DecodeString(publicKeyEncoded)
		if err != nil {
			return err
		}

		c.config.SigningPublicKey = new([32]byte)
		copy(c.config.SigningPublicKey[:], publicKeyDecoded)

		return nil
	}
}
