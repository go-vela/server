// SPDX-License-Identifier: Apache-2.0

package queue

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/go-vela/server/constants"
)

// Flags represents all supported command line
// interface (CLI) flags for the queue.
//
// https://pkg.go.dev/github.com/urfave/cli?tab=doc#Flag
var Flags = []cli.Flag{
	// Queue Flags

	&cli.StringFlag{
		Name:  "queue.driver",
		Usage: "driver to be used for the queue",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_QUEUE_DRIVER"),
			cli.EnvVar("QUEUE_DRIVER"),
			cli.File("/vela/queue/driver"),
		),
		Required: true,
	},
	&cli.StringFlag{
		Name:  "queue.addr",
		Usage: "fully qualified url (<scheme>://<host>) for the queue",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_QUEUE_ADDR"),
			cli.EnvVar("QUEUE_ADDR"),
			cli.File("/vela/queue/addr"),
		),
		Required: true,
		Action: func(_ context.Context, _ *cli.Command, v string) error {
			// check if the queue address has a scheme
			if !strings.Contains(v, "://") {
				return fmt.Errorf("queue address must be fully qualified (<scheme>://<host>)")
			}

			// check if the queue address has a trailing slash
			if strings.HasSuffix(v, "/") {
				return fmt.Errorf("queue address must not have trailing slash")
			}

			return nil
		},
	},
	&cli.BoolFlag{
		Name:  "queue.cluster",
		Usage: "enables connecting to a queue cluster",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_QUEUE_CLUSTER"),
			cli.EnvVar("QUEUE_CLUSTER"),
			cli.File("/vela/queue/cluster"),
		),
	},
	&cli.StringSliceFlag{
		Name:  "queue.routes",
		Usage: "list of routes (channels/topics) to publish builds",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_QUEUE_ROUTES"),
			cli.EnvVar("QUEUE_ROUTES"),
			cli.File("/vela/queue/routes"),
		),
		Value: []string{constants.DefaultRoute},
	},
	&cli.DurationFlag{
		Name:  "queue.pop.timeout",
		Usage: "timeout for requests that pop items off the queue",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_QUEUE_POP_TIMEOUT"),
			cli.EnvVar("QUEUE_POP_TIMEOUT"),
			cli.File("/vela/queue/pop_timeout"),
		),
		Value: 60 * time.Second,
	},
	&cli.StringFlag{
		Name:  "queue.private-key",
		Usage: "set value of base64 encoded queue signing private key",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_QUEUE_PRIVATE_KEY"),
			cli.EnvVar("QUEUE_PRIVATE_KEY"),
			cli.File("/vela/signing.key"),
		),
	},
	&cli.StringFlag{
		Name:  "queue.public-key",
		Usage: "set value of base64 encoded queue signing public key",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_QUEUE_PUBLIC_KEY"),
			cli.EnvVar("QUEUE_PUBLIC_KEY"),
			cli.File("/vela/signing.pub"),
		),
		Required: true,
	},
}
