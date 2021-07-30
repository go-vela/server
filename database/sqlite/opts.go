// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// ClientOpt represents a configuration option to initialize the database client.
type ClientOpt func(*client) error

// WithAddress sets the Sqlite address in the database client.
func WithAddress(address string) ClientOpt {
	logrus.Trace("configuring address in sqlite database client")

	return func(c *client) error {
		// check if the Sqlite address provided is empty
		if len(address) == 0 {
			return fmt.Errorf("no Sqlite address provided")
		}

		// set the address in the sqlite client
		c.config.Address = address

		return nil
	}
}

// WithCompressionLevel sets the Sqlite compression level in the database client.
func WithCompressionLevel(level int) ClientOpt {
	logrus.Trace("configuring compression level in sqlite database client")

	return func(c *client) error {
		// set the compression level in the sqlite client
		c.config.CompressionLevel = level

		return nil
	}
}

// WithConnectionLife sets the Sqlite connection duration in the database client.
func WithConnectionLife(duration time.Duration) ClientOpt {
	logrus.Trace("configuring connection duration in sqlite database client")

	return func(c *client) error {
		// set the connection duration in the sqlite client
		c.config.ConnectionLife = duration

		return nil
	}
}

// WithConnectionIdle sets the Sqlite maximum idle connections in the database client.
func WithConnectionIdle(idle int) ClientOpt {
	logrus.Trace("configuring maximum idle connections in sqlite database client")

	return func(c *client) error {
		// set the maximum idle connections in the sqlite client
		c.config.ConnectionIdle = idle

		return nil
	}
}

// WithConnectionOpen sets the Sqlite maximum open connections in the database client.
func WithConnectionOpen(open int) ClientOpt {
	logrus.Trace("configuring maximum open connections in sqlite database client")

	return func(c *client) error {
		// set the maximum open connections in the sqlite client
		c.config.ConnectionOpen = open

		return nil
	}
}

// WithEncryptionKey sets the Sqlite encryption key in the database client.
func WithEncryptionKey(key string) ClientOpt {
	logrus.Trace("configuring encryption key in sqlite database client")

	return func(c *client) error {
		// check if the Sqlite encryption key provided is empty
		if len(key) == 0 {
			return fmt.Errorf("no Sqlite encryption key provided")
		}

		// set the encryption key in the sqlite client
		c.config.EncryptionKey = key

		return nil
	}
}

// WithSkipCreation sets the Sqlite skip creation logic in the database client.
func WithSkipCreation(skipCreation bool) ClientOpt {
	logrus.Trace("configuring skip creating objects in sqlite database client")

	return func(c *client) error {
		// set the skip creating objects in the sqlite client
		c.config.SkipCreation = skipCreation

		return nil
	}
}
