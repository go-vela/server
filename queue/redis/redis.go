// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package redis

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type client struct {
	Queue   *redis.Client
	Options *redis.Options
}

// New returns a Queue implementation that
// integrates with a Redis queue instance.
func New(config string, channels ...string) (*client, error) {
	// parse the url provided
	options, err := redis.ParseURL(config)
	if err != nil {
		return nil, err
	}

	// create the Redis client from the parsed url
	queue := redis.NewClient(options)

	// setup queue with proper configuration
	err = setupQueue(queue)
	if err != nil {
		return nil, err
	}

	// create the client object
	client := &client{
		Queue:   queue,
		Options: options,
	}

	return client, nil
}

// NewCluster returns a Queue implementation that
// integrates with a Redis queue cluster.
func NewCluster(config string, channels ...string) (*client, error) {
	// parse the url provided
	options, err := redis.ParseURL(config)
	if err != nil {
		return nil, err
	}

	// create the Redis client from failover options
	queue := redis.NewFailoverClient(failoverFromOptions(options))

	// setup queue with proper configuration
	err = setupQueue(queue)
	if err != nil {
		return nil, err
	}

	// create the client object
	client := &client{
		Queue:   queue,
		Options: options,
	}

	return client, nil
}

// failoverFromOptions is a helper function to create
// the failover options from the parse options.
func failoverFromOptions(source *redis.Options) *redis.FailoverOptions {
	target := &redis.FailoverOptions{
		OnConnect:          source.OnConnect,
		Password:           source.Password,
		DB:                 source.DB,
		MaxRetries:         source.MaxRetries,
		MinRetryBackoff:    source.MinRetryBackoff,
		MaxRetryBackoff:    source.MaxRetryBackoff,
		DialTimeout:        source.DialTimeout,
		ReadTimeout:        source.ReadTimeout,
		WriteTimeout:       source.WriteTimeout,
		PoolSize:           source.PoolSize,
		MinIdleConns:       source.MinIdleConns,
		MaxConnAge:         source.MaxConnAge,
		PoolTimeout:        source.PoolTimeout,
		IdleTimeout:        source.IdleTimeout,
		IdleCheckFrequency: source.IdleCheckFrequency,
		TLSConfig:          source.TLSConfig,
	}

	// trim auto appended :6379 from address
	arrHosts := strings.TrimSuffix(source.Addr, ":6379")

	// remove array brackets from string
	// creating a comma separated list
	hosts := strings.TrimRight(
		strings.TrimLeft(arrHosts, "["), "]",
	)

	// the first host from the csv list is set as
	// the master node all subsequent hosts get
	// added as sentinel nodes
	for _, host := range strings.Split(hosts, ",") {
		if len(target.MasterName) == 0 {
			target.MasterName = host
			continue
		}

		target.SentinelAddrs = append(target.SentinelAddrs, host)
	}

	return target
}

// setupQueue is a helper function to setup the
// queue with the proper configuration.
func setupQueue(client *redis.Client) error {
	// ping the queue
	err := pingQueue(client)
	if err != nil {
		return err
	}

	return nil
}

// pingQueue is a helper function to send a "ping"
// request with backoff to the database.
//
// This will ensure we have properly established a
// connection to the Redis queue instance before
// we try to set it up.
func pingQueue(client *redis.Client) error {
	// attempt 10 times
	for i := 0; i < 10; i++ {
		// send ping request to queue
		err := client.Ping().Err()
		if err != nil {
			logrus.Debugf("Error pinging Redis queue. Retrying in %v", (time.Duration(i) * time.Second))
			time.Sleep(1 * time.Second)
			continue
		}

		return nil
	}

	return fmt.Errorf("Error establishing Redis queue connection")
}

// NewTest returns a Queue implementation that
// integrates with a local Redis instance.
//
// It's possible to overide this with env variables,
// which gets used as a part of integration testing
// with the different supported backends.
//
// This function is intended for running tests only.
func NewTest() (*client, error) {
	config := os.Getenv("VELA_QUEUE_CONFIG")
	if len(config) == 0 {
		config = "localhost:6379"
	}

	// parse the url provided
	options, err := redis.ParseURL(config)
	if err != nil {
		return nil, err
	}

	// create the Redis client from the parsed url
	queue := redis.NewClient(options)

	// create the client object
	client := &client{
		Queue:   queue,
		Options: options,
	}

	return client, nil
}
