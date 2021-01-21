// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/go-vela/server/source"
	"github.com/go-vela/server/source/github"

	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

// helper function to setup the source from the CLI arguments.
func setupSource(c *cli.Context) (source.Service, error) {
	logrus.Debug("Creating source client from CLI configuration")

	switch c.String("source-driver") {
	case constants.DriverGithub:
		return setupGithub(c)
	case constants.DriverGitlab:
		return setupGitlab(c)
	default:
		return nil, fmt.Errorf("invalid source driver: %s", c.String("source-driver"))
	}
}

// helper function to setup the GitHub source from the CLI arguments.
func setupGithub(c *cli.Context) (source.Service, error) {
	logrus.Tracef("Creating %s source client from CLI configuration", constants.DriverGithub)
	return github.New(c)
}

// helper function to setup the Gitlab source from the CLI arguments.
func setupGitlab(c *cli.Context) (source.Service, error) {
	logrus.Tracef("Creating %s source client from CLI configuration", constants.DriverGitlab)
	// return gitlab.New(c)
	return nil, fmt.Errorf("unsupported source driver: %s", constants.DriverGitlab)
}
