// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"fmt"
	"time"
)

// ClientOpt represents a configuration option to initialize the queue client for postgres.
type ClientOpt func(*client) error

// WithAddress sets the address in the queue client for Postgres.
func WithAddress(address string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring address in postgres queue client")

		// check if the address provided is empty
		if len(address) == 0 {
			return fmt.Errorf("no postgres queue address provided")
		}

		// set the queue address in the postgres client
		c.config.Address = address

		return nil
	}
}

// WithChannels sets the channels in the queue client for Postgres.
func WithChannels(channels ...string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring channels in postgres queue client")

		// check if the channels provided are empty
		if len(channels) == 0 {
			return fmt.Errorf("no postgres queue channels provided")
		}

		// set the queue channels in the postgres client
		c.config.Channels = channels

		return nil
	}
}

// WithPopTransactionTimeout sets the pop transaction timeout in the queue client for postgres.
func WithPopTransactionTimeout(timeout time.Duration) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring pop transaction timeout in postgres queue client")

		// set the queue pop transaction timeout in the postgres client
		c.config.PopTransactionTimeout = timeout

		return nil
	}
}

// WithPopTimeout sets the pop query timeout in the queue client for postgres.
func WithPopTimeout(timeout time.Duration) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring pop timeout in postgres queue client")

		// set the queue pop query timeout in the postgres client
		c.config.PopTimeout = timeout

		return nil
	}
}
