// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/go-vela/server/scm"
	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

// helper function to setup the source from the CLI arguments.
func setupSource(c *cli.Context) (scm.Service, error) {
	logrus.Debug("Creating source client from CLI configuration")

	// source configuration
	_setup := &scm.Setup{
		Driver:               c.String("scm.driver"),
		Address:              c.String("scm.addr"),
		ClientID:             c.String("scm.client"),
		ClientSecret:         c.String("scm.secret"),
		ServerAddress:        c.String("server-addr"),
		ServerWebhookAddress: c.String("scm.webhook.addr"),
		StatusContext:        c.String("scm.context"),
		WebUIAddress:         c.String("webui-addr"),
		Scopes:               c.StringSlice("scm.scopes"),
	}

	// setup the source
	//
	// https://pkg.go.dev/github.com/go-vela/server/scm?tab=doc#New
	return scm.New(_setup)
}
