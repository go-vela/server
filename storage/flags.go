// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	// STORAGE Flags

	&cli.BoolFlag{
		EnvVars:  []string{"VELA_STORAGE_ENABLE"},
		FilePath: "vela/storage/enable",
		Name:     "storage.enable",
		Usage:    "enable object storage",
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_STORAGE_DRIVER", "STORAGE_DRIVER"},
		FilePath: "vela/storage/driver",
		Name:     "storage.driver",
		Usage:    "object storage driver",
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_STORAGE_ADDRESS", "STORAGE_ADDRESS"},
		FilePath: "vela/storage/addr",
		Name:     "storage.addr",
		Usage:    "set the storage endpoint (ex. scheme://host:port)",
	},

	&cli.StringFlag{
		EnvVars:  []string{"VELA_STORAGE_ACCESS_KEY", "STORAGE_ACCESS_KEY"},
		FilePath: "vela/storage/access_key",
		Name:     "storage.access.key",
		Usage:    "set storage access key",
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_STORAGE_SECRET_KEY", "STORAGE_SECRET_KEY"},
		FilePath: "vela/storage/secret_key",
		Name:     "storage.secret.key",
		Usage:    "set storage secret key",
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_STORAGE_BUCKET"},
		FilePath: "vela/storage/bucket",
		Name:     "storage.bucket.name",
		Usage:    "set storage bucket name",
	},
	&cli.BoolFlag{
		EnvVars: []string{"VELA_STORAGE_USE_SSL"},
		Name:    "storage.use.ssl",
		Usage:   "enable storage to use SSL",
		Value:   false,
	},
}
