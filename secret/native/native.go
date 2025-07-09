// SPDX-License-Identifier: Apache-2.0

package native

import (
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database"
)

// Client represents a struct to hold native secret setup.
type Client struct {
	// client to interact with database for secret operations
	Database database.Interface
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
	Logger *logrus.Entry
}

// New returns a Secret implementation that integrates with a Native secrets engine.
func New(opts ...ClientOpt) (*Client, error) {
	// create new native client
	c := new(Client)

	// create new fields
	c.Database = *new(database.Interface)

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#StandardLogger
	logger := logrus.StandardLogger()

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#NewEntry
	c.Logger = logrus.NewEntry(logger).WithField("engine", c.Driver())

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}
