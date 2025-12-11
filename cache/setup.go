// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/cache/redis"
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

	InstallTokenKey string
}

// Redis creates and returns a Vela service capable
// of integrating with a Redis queue.
func (s *Setup) Redis(ctx context.Context) (Service, error) {
	logrus.Trace("creating redis queue client from setup")

	// create new Redis queue service
	//
	// https://pkg.go.dev/github.com/go-vela/server/queue/redis?tab=doc#New
	return redis.New(
		ctx,
		redis.WithAddress(s.Address),
		redis.WithCluster(s.Cluster),
		redis.WithInstallTokenKey(s.InstallTokenKey),
	)
}
