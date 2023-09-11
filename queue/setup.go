// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package queue

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-vela/server/queue/redis"
	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// Setup represents the configuration necessary for
// creating a Vela service capable of integrating
// with a configured queue environment.
type Setup struct {
	// Queue Configuration

	// specifies the driver to use for the queue client
	Driver string
	// specifies the address to use for the queue client
	Address string
	// enables the queue client to integrate with a queue cluster
	Cluster bool
	// specifies a list of routes (channels/topics) for managing builds for the queue client
	Routes []string
	// specifies the timeout for pop requests for the queue client
	Timeout time.Duration
	// private key in base64 used for signing items pushed to the queue
	PrivateKey string
	// public key in base64 used for opening items popped from the queue
	PublicKey string
}

// Redis creates and returns a Vela service capable
// of integrating with a Redis queue.
func (s *Setup) Redis() (Service, error) {
	logrus.Trace("creating redis queue client from setup")

	// create new Redis queue service
	//
	// https://pkg.go.dev/github.com/go-vela/server/queue/redis?tab=doc#New
	return redis.New(
		redis.WithAddress(s.Address),
		redis.WithChannels(s.Routes...),
		redis.WithCluster(s.Cluster),
		redis.WithTimeout(s.Timeout),
		redis.WithPrivateKey(s.PrivateKey),
		redis.WithPublicKey(s.PublicKey),
	)
}

// Kafka creates and returns a Vela service capable
// of integrating with a Kafka queue.
func (s *Setup) Kafka() (Service, error) {
	logrus.Trace("creating kafka queue client from setup")

	return nil, fmt.Errorf("unsupported queue driver: %s", constants.DriverKafka)
}

// Validate verifies the necessary fields for the
// provided configuration are populated correctly.
func (s *Setup) Validate() error {
	logrus.Trace("validating queue setup for client")

	// verify a queue driver was provided
	if len(s.Driver) == 0 {
		return fmt.Errorf("no queue driver provided")
	}

	// verify a queue address was provided
	if len(s.Address) == 0 {
		return fmt.Errorf("no queue address provided")
	}

	// check if the queue address has a scheme
	if !strings.Contains(s.Address, "://") {
		return fmt.Errorf("queue address must be fully qualified (<scheme>://<host>)")
	}

	// check if the queue address has a trailing slash
	if strings.HasSuffix(s.Address, "/") {
		return fmt.Errorf("queue address must not have trailing slash")
	}

	// verify queue routes were provided
	if len(s.Routes) == 0 {
		return fmt.Errorf("no queue routes provided")
	}

	if len(s.PublicKey) == 0 {
		return fmt.Errorf("no public key was provided")
	}

	if len(s.PrivateKey) == 0 {
		return fmt.Errorf("no private key was provided")
	}

	// setup is valid
	return nil
}
