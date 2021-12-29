// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"fmt"
	"time"
)

// ClientOpt represents a configuration option to initialize the database client for Postgres.
type ClientOpt func(*client) error

// WithAddress sets the address in the database client for Postgres.
func WithAddress(address string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring address in postgres database client")

		// check if the Postgres address provided is empty
		if len(address) == 0 {
			return fmt.Errorf("no Postgres address provided")
		}

		// set the address in the postgres client
		c.config.Address = address

		return nil
	}
}

// WithCompressionLevel sets the compression level in the database client for Postgres.
func WithCompressionLevel(level int) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring compression level in postgres database client")

		// set the compression level in the postgres client
		c.config.CompressionLevel = level

		return nil
	}
}

// WithConnectionLife sets the connection duration in the database client for Postgres.
func WithConnectionLife(duration time.Duration) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring connection duration in postgres database client")

		// set the connection duration in the postgres client
		c.config.ConnectionLife = duration

		return nil
	}
}

// WithConnectionIdle sets the maximum idle connections in the database client for Postgres.
func WithConnectionIdle(idle int) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring maximum idle connections in postgres database client")

		// set the maximum idle connections in the postgres client
		c.config.ConnectionIdle = idle

		return nil
	}
}

// WithConnectionOpen sets the maximum open connections in the database client for Postgres.
func WithConnectionOpen(open int) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring maximum open connections in postgres database client")

		// set the maximum open connections in the postgres client
		c.config.ConnectionOpen = open

		return nil
	}
}

// WithEncryptionKey sets the encryption key in the database client for Postgres.
func WithEncryptionKey(key string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring encryption key in postgres database client")

		// check if the Postgres encryption key provided is empty
		if len(key) == 0 {
			return fmt.Errorf("no Postgres encryption key provided")
		}

		// set the encryption key in the postgres client
		c.config.EncryptionKey = key

		return nil
	}
}

// WithSkipCreation sets the skip creation logic in the database client for Postgres.
func WithSkipCreation(skipCreation bool) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring skip creating objects in postgres database client")

		// set to skip creating tables and indexes in the postgres client
		c.config.SkipCreation = skipCreation

		return nil
	}
}
