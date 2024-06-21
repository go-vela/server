// SPDX-License-Identifier: Apache-2.0

package queue

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/go-vela/types/constants"
)

// FromCLIContext helper function to setup the queue from the CLI arguments.
func FromCLIContext(c *cli.Context) (Service, error) {
	logrus.Debug("creating queue client from CLI configuration")

	// queue configuration
	_setup := &Setup{
		Driver:     c.String("queue.driver"),
		Address:    c.String("queue.addr"),
		Cluster:    c.Bool("queue.cluster"),
		Routes:     c.StringSlice("queue.routes"),
		Timeout:    c.Duration("queue.pop.timeout"),
		PrivateKey: c.String("queue.private-key"),
		PublicKey:  c.String("queue.public-key"),
	}

	// setup the queue
	//
	// https://pkg.go.dev/github.com/go-vela/server/queue?tab=doc#New
	return New(_setup)
}

// New creates and returns a Vela service capable of
// integrating with the configured queue environment.
// Currently, the following queues are supported:
//
// * redis
// .
func New(s *Setup) (Service, error) {
	// validate the setup being provided
	//
	// https://pkg.go.dev/github.com/go-vela/server/queue?tab=doc#Setup.Validate
	err := s.Validate()
	if err != nil {
		return nil, err
	}

	logrus.Debug("creating queue client from setup")
	// process the queue driver being provided
	switch s.Driver {
	case constants.DriverKafka:
		// handle the Kafka queue driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/server/queue?tab=doc#Setup.Kafka
		return s.Kafka()
	case constants.DriverRedis:
		// handle the Redis queue driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/server/queue?tab=doc#Setup.Redis
		return s.Redis()
	default:
		// handle an invalid queue driver being provided
		return nil, fmt.Errorf("invalid queue driver provided: %s", s.Driver)
	}
}
