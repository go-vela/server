// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	// STORAGE Flags

	&cli.BoolFlag{
		EnvVars: []string{"VELA_STORAGE_ENABLE"},
		Name:    "storage.enable",
		Usage:   "enable object storage",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_STORAGE_DRIVER"},
		Name:    "storage.driver.name",
		Usage:   "object storage driver",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_STORAGE_ENDPOINT"},
		Name:    "storage.endpoint.name",
		Usage:   "set the storage endpoint (ex. scheme://host:port)",
	},

	&cli.StringFlag{
		EnvVars: []string{"VELA_STORAGE_ACCESS_KEY"},
		Name:    "storage.access.key",
		Usage:   "set storage access key",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_STORAGE_SECRET_KEY"},
		Name:    "storage.secret.key",
		Usage:   "set storage secret key",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_STORAGE_BUCKET"},
		Name:    "storage.bucket.name",
		Usage:   "set storage bucket name",
	},
	&cli.BoolFlag{
		EnvVars: []string{"VELA_STORAGE_USE_SSL"},
		Name:    "storage.use.ssl",
		Usage:   "enable storage to use SSL",
		Value:   false,
	},
}
