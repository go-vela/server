// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"fmt"
	"time"
)

// ClientOpt represents a configuration option to initialize the database client for Sqlite.
type ClientOpt func(*client) error

// WithAddress sets the Sqlite address in the database client for Sqlite.
func WithAddress(address string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring address in sqlite database client")

		// check if the Sqlite address provided is empty
		if len(address) == 0 {
			return fmt.Errorf("no Sqlite address provided")
		}

		// set the address in the sqlite client
		c.config.Address = address

		return nil
	}
}

// WithCompressionLevel sets the compression level in the database client for Sqlite.
func WithCompressionLevel(level int) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring compression level in sqlite database client")

		// set the compression level in the sqlite client
		c.config.CompressionLevel = level

		return nil
	}
}

// WithConnectionLife sets the connection duration in the database client for Sqlite.
func WithConnectionLife(duration time.Duration) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring connection duration in sqlite database client")

		// set the connection duration in the sqlite client
		c.config.ConnectionLife = duration

		return nil
	}
}

// WithConnectionIdle sets the maximum idle connections in the database client for Sqlite.
func WithConnectionIdle(idle int) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring maximum idle connections in sqlite database client")

		// set the maximum idle connections in the sqlite client
		c.config.ConnectionIdle = idle

		return nil
	}
}

// WithConnectionOpen sets the maximum open connections in the database client for Sqlite.
func WithConnectionOpen(open int) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring maximum open connections in sqlite database client")

		// set the maximum open connections in the sqlite client
		c.config.ConnectionOpen = open

		return nil
	}
}

// WithEncryptionKey sets the encryption key in the database client for Sqlite.
func WithEncryptionKey(key string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring encryption key in sqlite database client")

		// check if the Sqlite encryption key provided is empty
		if len(key) == 0 {
			return fmt.Errorf("no Sqlite encryption key provided")
		}

		// set the encryption key in the sqlite client
		c.config.EncryptionKey = key

		return nil
	}
}

// WithSkipCreation sets the skip creation logic in the database client for Sqlite.
func WithSkipCreation(skipCreation bool) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring skip creating objects in sqlite database client")

		// set to skip creating tables and indexes in the sqlite client
		c.config.SkipCreation = skipCreation

		return nil
	}
}
