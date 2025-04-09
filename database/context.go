// SPDX-License-Identifier: Apache-2.0

package database

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/go-vela/server/tracing"
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

	return d
}

// ToContext adds the database Interface to this context if it supports
// the Setter interface.
func ToContext(c Setter, d Interface) {
	c.Set(key, d)
}

// FromCLICommand creates and returns a database engine from the urfave/cli context.
func FromCLICommand(c *cli.Command, tc *tracing.Client) (Interface, error) {
	logrus.Debug("creating database engine from CLI configuration")

	return New(
		WithAddress(c.String("database.addr")),
		WithCompressionLevel(int(c.Int("database.compression.level"))),
		WithConnectionLife(c.Duration("database.connection.life")),
		WithConnectionIdle(int(c.Int("database.connection.idle"))),
		WithConnectionOpen(int(c.Int("database.connection.open"))),
		WithDriver(c.String("database.driver")),
		WithEncryptionKey(c.String("database.encryption.key")),
		WithLogLevel(c.String("database.log.level")),
		WithLogSkipNotFound(c.Bool("database.log.skip_notfound")),
		WithLogSlowThreshold(c.Duration("database.log.slow_threshold")),
		WithLogShowSQL(c.Bool("database.log.show_sql")),
		WithSkipCreation(c.Bool("database.skip_creation")),
		WithTracing(tc),
	)
}
