// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/go-vela/server/queue"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/urfave/cli/v2"
)

// helper function to setup the queue from the CLI arguments.
func setupQueue(c *cli.Context, s *api.Settings) (queue.Service, error) {
	logrus.Debug("Creating queue client from CLI configuration")

	// queue configuration
	_setup := &queue.Setup{
		Driver:     c.String("queue.driver"),
		Address:    c.String("queue.addr"),
		Cluster:    c.Bool("queue.cluster"),
		Routes:     s.GetQueueRoutes(),
		Timeout:    c.Duration("queue.pop.timeout"),
		PrivateKey: c.String("queue.private-key"),
		PublicKey:  c.String("queue.public-key"),
	}

	// setup the queue
	//
	// https://pkg.go.dev/github.com/go-vela/server/queue?tab=doc#New
	q, err := queue.New(_setup)
	if err != nil {
		return nil, err
	}

	return q, nil
}
