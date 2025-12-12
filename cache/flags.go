// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
)

// Flags represents all supported command line
// interface (CLI) flags for the queue.
//
// https://pkg.go.dev/github.com/urfave/cli?tab=doc#Flag
var Flags = []cli.Flag{
	// Queue Flags

	&cli.StringFlag{
		Name:  "cache.driver",
		Usage: "driver to be used for the cache",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_CACHE_DRIVER"),
			cli.EnvVar("CACHE_DRIVER"),
			cli.File("/vela/cache/driver"),
		),
		Required: true,
	},
	&cli.StringFlag{
		Name:  "cache.addr",
		Usage: "fully qualified url (<scheme>://<host>) for the cache",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_CACHE_ADDR"),
			cli.EnvVar("CACHE_ADDR"),
			cli.File("/vela/cache/addr"),
		),
		Action: func(_ context.Context, _ *cli.Command, v string) error {
			// check if the cache address has a scheme
			if !strings.Contains(v, "://") {
				return fmt.Errorf("cache address must be fully qualified (<scheme>://<host>)")
			}

			// check if the cache address has a trailing slash
			if strings.HasSuffix(v, "/") {
				return fmt.Errorf("cache address must not have trailing slash")
			}

			return nil
		},
	},
	&cli.BoolFlag{
		Name:  "cache.cluster",
		Usage: "enables connecting to a cache cluster",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_CACHE_CLUSTER"),
			cli.EnvVar("CACHE_CLUSTER"),
			cli.File("/vela/cache/cluster"),
		),
	},
	&cli.StringFlag{
		Name:  "cache.install-token-key",
		Usage: "set cache install token key",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_CACHE_INSTALL_TOKEN_KEY"),
			cli.EnvVar("CACHE_INSTALL_TOKEN_KEY"),
			cli.File("/vela/cache/install_token_key"),
		),
		Required: true,
	},
}
