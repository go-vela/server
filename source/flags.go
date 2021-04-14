// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package source

import (
	"github.com/go-vela/types/constants"
	"github.com/urfave/cli/v2"
)

// Flags represents all supported command line
// interface (CLI) flags for the source.
//
// https://pkg.go.dev/github.com/urfave/cli?tab=doc#Flag
var Flags = []cli.Flag{

	// Logger Flags

	&cli.StringFlag{
		EnvVars: []string{"VELA_SOURCE_LOG_FORMAT", "SOURCE_LOG_FORMAT", "VELA_LOG_FORMAT"},
		Name:    "source.log.format",
		Usage:   "format of logs to output",
		Value:   "json",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_SOURCE_LOG_LEVEL", "SOURCE_LOG_LEVEL", "VELA_LOG_LEVEL"},
		Name:    "source.log.level",
		Usage:   "level of logs to output",
		Value:   "info",
	},

	// Source Flags

	&cli.StringFlag{
		EnvVars: []string{"VELA_SOURCE_DRIVER", "SOURCE_DRIVER"},
		Name:    "source.driver",
		Usage:   "driver to be used for the version control system",
		Value:   constants.DriverGithub,
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_SOURCE_ADDR", "SOURCE_ADDR"},
		Name:    "source.addr",
		Usage:   "fully qualified url (<scheme>://<host>) for the version control system",
		Value:   "https://github.com",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_SOURCE_CLIENT", "SOURCE_CLIENT"},
		Name:    "source.client",
		Usage:   "OAuth client id from version control system",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_SOURCE_SECRET", "SOURCE_SECRET"},
		Name:    "source.secret",
		Usage:   "OAuth client secret from version control system",
	},
	&cli.StringFlag{
		EnvVars: []string{"VELA_SOURCE_CONTEXT", "SOURCE_CONTEXT"},
		Name:    "source.context",
		Usage:   "context for commit status in version control system",
		Value:   "continuous-integration/vela",
	},
}
