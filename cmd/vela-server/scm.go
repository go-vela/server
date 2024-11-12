// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/tracing"
)

// helper function to setup the scm from the CLI arguments.
func setupSCM(c *cli.Context, tc *tracing.Client) (scm.Service, error) {
	logrus.Debug("creating scm client from CLI configuration")

	// scm configuration
	_setup := &scm.Setup{
		Driver:               c.String("scm.driver"),
		Address:              c.String("scm.addr"),
		ClientID:             c.String("scm.client"),
		ClientSecret:         c.String("scm.secret"),
		AppID:                c.Int64("scm.app.id"),
		AppPrivateKey:        c.String("scm.app.private_key"),
		AppPermissions:       c.StringSlice("scm.app.permissions"),
		ServerAddress:        c.String("server-addr"),
		ServerWebhookAddress: c.String("scm.webhook.addr"),
		StatusContext:        c.String("scm.context"),
		WebUIAddress:         c.String("webui-addr"),
		OAuthScopes:          c.StringSlice("scm.scopes"),
		Tracing:              tc,
	}

	// setup the scm
	//
	// https://pkg.go.dev/github.com/go-vela/server/scm?tab=doc#New
	return scm.New(c.Context, _setup)
}
