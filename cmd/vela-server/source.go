// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/go-vela/server/source"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

// helper function to setup the source from the CLI arguments.
func setupSource(c *cli.Context) (source.Service, error) {
	logrus.Debug("Creating source client from CLI configuration")

	// source configuration
	_setup := &source.Setup{
		Driver:        c.String("source.driver"),
		Address:       c.String("source.addr"),
		ClientID:      c.String("source.client"),
		ClientSecret:  c.String("source.secret"),
		ServerAddress: c.String("server-addr"),
		ServerWebhookAddress: c.String("source.webhook.addr"),
		StatusContext: c.String("source.context"),
		WebUIAddress:  c.String("webui-addr"),
		Scopes:        c.StringSlice("source.scopes"),
	}

	// setup the source
	//
	// https://pkg.go.dev/github.com/go-vela/server/source?tab=doc#New
	return source.New(_setup)
}
