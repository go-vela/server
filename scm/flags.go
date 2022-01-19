// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package scm

import (
	"github.com/go-vela/types/constants"
	"github.com/urfave/cli/v2"
)

// Flags represents all supported command line
// interface (CLI) flags for the scm.
//
// https://pkg.go.dev/github.com/urfave/cli?tab=doc#Flag
//
// TODO: in a future release remove the "source" vars in favor of the "scm" ones.
// nolint:lll // these errors will go away when the TODO is completed
var Flags = []cli.Flag{
	// SCM Flags

	&cli.StringFlag{
		EnvVars:  []string{"VELA_SCM_DRIVER", "SCM_DRIVER", "VELA_SOURCE_DRIVER", "SOURCE_DRIVER"},
		FilePath: "/vela/scm/driver",
		Name:     "scm.driver",
		Usage:    "driver to be used for the version control system",
		Value:    constants.DriverGithub,
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_SCM_ADDR", "SCM_ADDR", "VELA_SOURCE_ADDR", "SOURCE_ADDR"},
		FilePath: "/vela/scm/addr",
		Name:     "scm.addr",
		Usage:    "fully qualified url (<scheme>://<host>) for the version control system",
		Value:    "https://github.com",
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_SCM_CLIENT", "SCM_CLIENT", "VELA_SOURCE_CLIENT", "SOURCE_CLIENT"},
		FilePath: "/vela/scm/client",
		Name:     "scm.client",
		Usage:    "OAuth client id from version control system",
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_SCM_SECRET", "SCM_SECRET", "VELA_SOURCE_SECRET", "SOURCE_SECRET"},
		FilePath: "/vela/scm/secret",
		Name:     "scm.secret",
		Usage:    "OAuth client secret from version control system",
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_SCM_CONTEXT", "SCM_CONTEXT", "VELA_SOURCE_CONTEXT", "SOURCE_CONTEXT"},
		FilePath: "/vela/scm/context",
		Name:     "scm.context",
		Usage:    "context for commit status in version control system",
		Value:    "continuous-integration/vela",
	},
	&cli.StringSliceFlag{
		EnvVars:  []string{"VELA_SCM_SCOPES", "SCM_SCOPES", "VELA_SOURCE_SCOPES", "SOURCE_SCOPES"},
		FilePath: "/vela/scm/scopes",
		Name:     "scm.scopes",
		Usage:    "OAuth scopes to be used for the version control system",
		Value:    cli.NewStringSlice("repo", "repo:status", "user:email", "read:user", "read:org"),
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_SCM_WEBHOOK_ADDR", "SCM_WEBHOOK_ADDR", "VELA_SOURCE_WEBHOOK_ADDR", "SOURCE_WEBHOOK_ADDR"},
		FilePath: "/vela/scm/webhook_addr",
		Name:     "scm.webhook.addr",
		Usage: "Alternative or proxy server address as a fully qualified url (<scheme>://<host>). " +
			"Use this when the Vela server address that the scm provider can send webhooks to " +
			"differs from the server address the UI and oauth flows use, such as when the server " +
			"is behind a Firewall or NAT, or when using something like ngrok to forward webhooks. " +
			"(defaults to VELA_ADDR).",
	},
}
