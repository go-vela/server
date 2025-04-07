// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"strings"
)

var Flags = []cli.Flag{
	// STORAGE Flags

	&cli.BoolFlag{
		Name:  "storage.enable",
		Usage: "enable object storage",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_STORAGE_ENABLE"),
			cli.File("vela/storage/enable"),
		),
		Required: true,
	},
	&cli.StringFlag{
		Name:  "storage.driver",
		Usage: "object storage driver",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_STORAGE_DRIVER"),
			cli.EnvVar("STORAGE_DRIVER"),
			cli.File("vela/storage/driver"),
		),
		Required: true,
	},
	&cli.StringFlag{
		Name:  "storage.addr",
		Usage: "set the storage endpoint (ex. scheme://host:port)",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_STORAGE_ADDRESS"),
			cli.EnvVar("STORAGE_ADDRESS"),
			cli.File("vela/storage/addr"),
		),
		Required: true,
		Action: func(_ context.Context, _ *cli.Command, v string) error {
			// check if the queue address has a scheme
			if !strings.Contains(v, "://") {
				return fmt.Errorf("storage address must be fully qualified (<scheme>://<host>)")
			}

			// check if the queue address has a trailing slash
			if strings.HasSuffix(v, "/") {
				return fmt.Errorf("storage address must not have trailing slash")
			}

			return nil
		},
	},

	&cli.StringFlag{
		Name:  "storage.access.key",
		Usage: "set storage access key",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_STORAGE_ACCESS_KEY"),
			cli.EnvVar("STORAGE_ACCESS_KEY"),
			cli.File("vela/storage/access_key"),
		),
		Required: true,
	},
	&cli.StringFlag{
		Name:  "storage.secret.key",
		Usage: "set storage secret key",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_STORAGE_SECRET_KEY"),
			cli.EnvVar("STORAGE_SECRET_KEY"),
			cli.File("vela/storage/secret_key"),
		),
		Required: true,
	},
	&cli.StringFlag{
		Name:  "storage.bucket.name",
		Usage: "set storage bucket name",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_STORAGE_BUCKET"),
			cli.File("vela/storage/bucket"),
		),
		Required: true,
	},
	&cli.BoolFlag{
		Name:  "storage.use.ssl",
		Usage: "enable storage to use SSL",
		Value: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_STORAGE_USE_SSL"),
		),
		Required: false,
	},
}
