// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/tracing"
)

// helper function to setup the scm from the CLI arguments.
func setupSCM(ctx context.Context, c *cli.Command, tc *tracing.Client) (scm.Service, error) {
	logrus.Debug("creating scm client from CLI configuration")

	// scm configuration
	_setup := &scm.Setup{
		Driver:               c.String("scm.driver"),
		Address:              c.String("scm.addr"),
		ClientID:             c.String("scm.client"),
		ClientSecret:         c.String("scm.secret"),
		AppID:                c.Int("scm.app.id"),
		AppPrivateKey:        c.String("scm.app.private-key"),
		AppPrivateKeyPath:    c.String("scm.app.private-key.path"),
		AppPermissions:       c.StringSlice("scm.app.permissions"),
		ServerAddress:        c.String("server-addr"),
		ServerWebhookAddress: c.String("scm.webhook.addr"),
		StatusContext:        c.String("scm.context"),
		WebUIAddress:         c.String("webui-addr"),
		OAuthScopes:          c.StringSlice("scm.scopes"),
		RepoRoleMap:          c.StringMap("scm.repo.roles-map"),
		OrgRoleMap:           c.StringMap("scm.org.roles-map"),
		TeamRoleMap:          c.StringMap("scm.team.roles-map"),
		Tracing:              tc,
	}

	// setup the scm
	//
	// https://pkg.go.dev/github.com/go-vela/server/scm?tab=doc#New
	return scm.New(ctx, _setup)
}
