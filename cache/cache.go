// SPDX-License-Identifier: Apache-2.0
package cache

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/go-vela/server/constants"
)

// FromCLICommand helper function to setup the queue from the CLI arguments.
func FromCLICommand(ctx context.Context, c *cli.Command) (Service, error) {
	logrus.Debug("creating queue client from CLI configuration")

	// queue configuration
	_setup := &Setup{
		Driver:          c.String("cache.driver"),
		Address:         c.String("cache.addr"),
		Cluster:         c.Bool("cache.cluster"),
		InstallTokenKey: c.String("cache.install-token-key"),
	}

	// setup the queue
	//
	// https://pkg.go.dev/github.com/go-vela/server/queue?tab=doc#New
	return New(ctx, _setup)
}

// New creates and returns a Vela service capable of
// integrating with the configured queue environment.
// Currently, the following queues are supported:
//
// * redis
// .
func New(ctx context.Context, s *Setup) (Service, error) {
	logrus.Debug("creating queue client from setup")
	// process the queue driver being provided
	switch s.Driver {
	case constants.DriverRedis:
		// handle the Redis queue driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/server/queue?tab=doc#Setup.Redis
		return s.Redis(ctx)
	default:
		// handle an invalid queue driver being provided
		return nil, fmt.Errorf("invalid queue driver provided: %s", s.Driver)
	}
}
