// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/queue/redis"

	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

// helper function to setup the queue from the CLI arguments.
func setupQueue(c *cli.Context) (queue.Service, error) {
	logrus.Debug("Creating queue client from CLI configuration")
	switch c.String("queue-driver") {
	case constants.DriverKafka:
		return setupKafka(c)
	case constants.DriverRedis:
		return setupRedis(c)
	default:
		return nil, fmt.Errorf("Unrecognized queue driver: %s", c.String("queue-driver"))
	}
}

// helper function to setup the Kafka queue from the CLI arguments.
func setupKafka(c *cli.Context) (queue.Service, error) {
	logrus.Tracef("Creating %s queue client from CLI configuration", constants.DriverKafka)
	// return kafka.New(c.String("queue-config"), "vela")
	return nil, fmt.Errorf("Unsupported queue driver: %s", constants.DriverKafka)
}

// helper function to setup the Redis queue from the CLI arguments.
func setupRedis(c *cli.Context) (queue.Service, error) {

	// setup routes
	routes := append(c.StringSlice("queue-worker-routes"), constants.DefaultRoute)

	if c.Bool("queue-cluster") {
		logrus.Tracef("Creating %s queue cluster client from CLI configuration", constants.DriverRedis)
		return redis.NewCluster(c.String("queue-config"), routes...)
	}

	logrus.Tracef("Creating %s queue client from CLI configuration", constants.DriverRedis)
	return redis.New(c.String("queue-config"), routes...)
}
