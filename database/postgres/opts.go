// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// ClientOpt represents a configuration option to initialize the database client.
type ClientOpt func(*client) error

// WithAddress sets the Postgres address in the database client.
func WithAddress(address string) ClientOpt {
	logrus.Trace("configuring address in postgres database client")

	return func(c *client) error {
		// check if the Postgres address provided is empty
		if len(address) == 0 {
			return fmt.Errorf("no Postgres address provided")
		}

		// set the address in the postgres client
		c.config.Address = address

		return nil
	}
}

// WithCompressionLevel sets the Postgres compression level in the database client.
func WithCompressionLevel(level int) ClientOpt {
	logrus.Trace("configuring compression level in postgres database client")

	return func(c *client) error {
		// set the compression level in the postgres client
		c.config.CompressionLevel = level

		return nil
	}
}

// WithConnectionLife sets the Postgres connection duration in the database client.
func WithConnectionLife(duration time.Duration) ClientOpt {
	logrus.Trace("configuring connection duration in postgres database client")

	return func(c *client) error {
		// set the connection duration in the postgres client
		c.config.ConnectionLife = duration

		return nil
	}
}

// WithConnectionIdle sets the Postgres maximum idle connections in the database client.
func WithConnectionIdle(idle int) ClientOpt {
	logrus.Trace("configuring maximum idle connections in postgres database client")

	return func(c *client) error {
		// set the maximum idle connections in the postgres client
		c.config.ConnectionIdle = idle

		return nil
	}
}

// WithConnectionOpen sets the Postgres maximum open connections in the database client.
func WithConnectionOpen(open int) ClientOpt {
	logrus.Trace("configuring maximum open connections in postgres database client")

	return func(c *client) error {
		// set the maximum open connections in the postgres client
		c.config.ConnectionOpen = open

		return nil
	}
}

// WithEncryptionKey sets the Postgres encryption key in the database client.
func WithEncryptionKey(key string) ClientOpt {
	logrus.Trace("configuring encryption key in postgres database client")

	return func(c *client) error {
		// check if the Postgres encryption key provided is empty
		if len(key) == 0 {
			return fmt.Errorf("no Postgres encryption key provided")
		}

		// set the encryption key in the postgres client
		c.config.EncryptionKey = key

		return nil
	}
}

// WithSkipCreation sets the Postgres skip creation logic in the database client.
func WithSkipCreation(skipCreation bool) ClientOpt {
	logrus.Trace("configuring skip creating objects in postgres database client")

	return func(c *client) error {
		// set to skip creating tables and indexes in the postgres client
		c.config.SkipCreation = skipCreation

		return nil
	}
}
