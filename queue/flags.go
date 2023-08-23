// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package queue

import (
	"time"

	"github.com/go-vela/types/constants"
	"github.com/urfave/cli/v2"
)

// Flags represents all supported command line
// interface (CLI) flags for the queue.
//
// https://pkg.go.dev/github.com/urfave/cli?tab=doc#Flag
var Flags = []cli.Flag{
	// Queue Flags

	&cli.StringFlag{
		EnvVars:  []string{"VELA_QUEUE_DRIVER", "QUEUE_DRIVER"},
		FilePath: "/vela/queue/driver",
		Name:     "queue.driver",
		Usage:    "driver to be used for the queue",
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_QUEUE_ADDR", "QUEUE_ADDR"},
		FilePath: "/vela/queue/addr",
		Name:     "queue.addr",
		Usage:    "fully qualified url (<scheme>://<host>) for the queue",
	},
	&cli.BoolFlag{
		EnvVars:  []string{"VELA_QUEUE_CLUSTER", "QUEUE_CLUSTER"},
		FilePath: "/vela/queue/cluster",
		Name:     "queue.cluster",
		Usage:    "enables connecting to a queue cluster",
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"VELA_QUEUE_ROUTES", "QUEUE_ROUTES"},
		FilePath: "/vela/queue/routes",
		Name:     "queue.routes",
		Usage:    "list of routes (channels/topics) to publish builds",
		Value:    cli.NewStringSlice(constants.DefaultRoute),
	},
	&cli.DurationFlag{
		EnvVars:  []string{"VELA_QUEUE_POP_TIMEOUT", "QUEUE_POP_TIMEOUT"},
		FilePath: "/vela/queue/pop_timeout",
		Name:     "queue.pop.timeout",
		Usage:    "timeout for requests that pop items off the queue",
		Value:    60 * time.Second,
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_QUEUE_SIGNING_PRIVATE_KEY"},
		FilePath: "/vela/signing.key",
		Name:     "queue.private-key",
		Usage:    "set value of base64 encoded queue signing private key",
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_QUEUE_SIGNING_PUBLIC_KEY"},
		FilePath: "/vela/signing.pub",
		Name:     "queue.public-key",
		Usage:    "set value of base64 encoded queue signing public key",
	},
}
