// SPDX-License-Identifier: Apache-2.0

package queue

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/constants"
)

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
