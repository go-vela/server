// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package redis

import (
	"context"
	"fmt"
	"strings"

	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type config struct {
	// specifies the address to use for the Redis client
	Address string
	// specifies a list of channels for managing builds for the Redis client
	Channels []string
	// enables the Redis client to integrate with a Redis cluster
	Cluster bool
	// specifies the timeout to use for the Redis client
	Timeout time.Duration
}

type client struct {
	config  *config
	Queue   *redis.Client
	Options *redis.Options
}

// New returns a Queue implementation that
// integrates with a Redis queue instance.
//
// nolint: golint // ignore returning unexported client
func New(opts ...ClientOpt) (*client, error) {
	// create new Redis client
	c := new(client)

	// create new fields
	c.config = new(config)
	c.Queue = new(redis.Client)
	c.Options = new(redis.Options)

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// parse the url provided
	options, err := redis.ParseURL(c.config.Address)
	if err != nil {
		return nil, err
	}

	// create the Redis options from the parsed url
	c.Options = options

	// check if clustering mode is enabled
	if c.config.Cluster {
		// create the Redis cluster client from the options
		c.Queue = redis.NewFailoverClient(failoverFromOptions(c.Options))
	} else {
		// create the Redis client from the parsed url
		c.Queue = redis.NewClient(c.Options)
	}

	// ping the queue
	err = pingQueue(c.Queue)
	if err != nil {
		return nil, err
	}

	return c, nil
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

// pingQueue is a helper function to send a "ping"
// request with backoff to the database.
//
// This will ensure we have properly established a
// connection to the Redis queue instance before
// we try to set it up.
func pingQueue(client *redis.Client) error {
	// attempt 10 times
	for i := 0; i < 10; i++ {
		// send ping request to client
		err := client.Ping(context.Background()).Err()
		if err != nil {
			logrus.Debugf("unable to ping Redis queue. Retrying in %v", (time.Duration(i) * time.Second))
			time.Sleep(1 * time.Second)

			continue
		}

		return nil
	}

	return fmt.Errorf("unable to establish connection to Redis queue")
}

// NewTest returns a Queue implementation that
// integrates with a local Redis instance.
//
// It's possible to overide this with env variables,
// which gets used as a part of integration testing
// with the different supported backends.
//
// This function is intended for running tests only.
//
// nolint: golint // ignore returning unexported client
func NewTest(channels ...string) (*client, error) {
	// create a local fake redis instance
	//
	// https://pkg.go.dev/github.com/alicebob/miniredis/v2#Run
	_redis, err := miniredis.Run()
	if err != nil {
		return nil, err
	}

	return New(
		WithAddress(fmt.Sprintf("redis://%s", _redis.Addr())),
		WithChannels(channels...),
		WithCluster(false),
	)
}
