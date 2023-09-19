// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"context"

	"github.com/go-vela/server/tracing"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const key = "database"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the database Interface associated with this context.
func FromContext(c context.Context) Interface {
	v := c.Value(key)
	if v == nil {
		return nil
	}

	d, ok := v.(Interface)
	if !ok {
		return nil
	}

	// assign as inflight property

	return d
}

// ToContext adds the database Interface to this context if it supports
// the Setter interface.
func ToContext(c Setter, d Interface) {
	c.Set(key, d)
}

// FromCLIContext creates and returns a database engine from the urfave/cli context.
func FromCLIContext(c *cli.Context, tc *tracing.Config) (Interface, error) {
	logrus.Debug("creating database engine from CLI configuration")

	return New(
		WithAddress(c.String("database.addr")),
		WithCompressionLevel(c.Int("database.compression.level")),
		WithConnectionLife(c.Duration("database.connection.life")),
		WithConnectionIdle(c.Int("database.connection.idle")),
		WithConnectionOpen(c.Int("database.connection.open")),
		WithDriver(c.String("database.driver")),
		WithEncryptionKey(c.String("database.encryption.key")),
		WithSkipCreation(c.Bool("database.skip_creation")),
		WithTracingConfig(tc),
	)
}
