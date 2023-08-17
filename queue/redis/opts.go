// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package redis

import (
	"encoding/base64"
	"errors"
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
//
//nolint:dupl // ignore similar code
func WithPrivateKey(privateKeyEncoded string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring private key in redis queue client")

		if len(privateKeyEncoded) == 0 {
			return errors.New("unable to base64 decode private key, provided key is empty")
		}

		privateKeyDecoded, err := base64.StdEncoding.DecodeString(privateKeyEncoded)
		if err != nil {
			return err
		}

		if len(privateKeyDecoded) == 0 {
			return errors.New("unable to base64 decode private key, decoded key is empty")
		}

		c.config.SigningPrivateKey = new([64]byte)
		copy(c.config.SigningPrivateKey[:], privateKeyDecoded)

		if c.config.SigningPrivateKey == nil {
			return errors.New("unable to copy decoded queue signing private key, copied key is nil")
		}

		if len(c.config.SigningPrivateKey) == 0 {
			return errors.New("unable to copy decoded queue signing private key, copied key is empty")
		}

		return nil
	}
}

// WithPublicKey sets the public key in the queue client for Redis.
//
//nolint:dupl // ignore similar code
func WithPublicKey(publicKeyEncoded string) ClientOpt {
	return func(c *client) error {
		c.Logger.Tracef("configuring public key in redis queue client")

		if len(publicKeyEncoded) == 0 {
			return errors.New("unable to base64 decode public key, provided key is empty")
		}

		publicKeyDecoded, err := base64.StdEncoding.DecodeString(publicKeyEncoded)
		if err != nil {
			return err
		}

		if len(publicKeyDecoded) == 0 {
			return errors.New("unable to base64 decode public key, decoded key is empty")
		}

		c.config.SigningPublicKey = new([32]byte)
		copy(c.config.SigningPublicKey[:], publicKeyDecoded)

		if c.config.SigningPublicKey == nil {
			return errors.New("unable to copy decoded queue signing public key, copied key is nil")
		}

		if len(c.config.SigningPublicKey) == 0 {
			return errors.New("unable to copy decoded queue signing public key, copied key is empty")
		}

		return nil
	}
}
