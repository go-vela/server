// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/go-vela/compiler/compiler"
	"github.com/go-vela/compiler/compiler/native"

	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

// helper function to setup the queue from the CLI arguments.
func setupCompiler(c *cli.Context) (compiler.Engine, error) {
	logrus.Debug("Creating queue client from CLI configuration")
	return setupCompilerNative(c)
}

// helper function to setup the Kafka queue from the CLI arguments.
func setupCompilerNative(c *cli.Context) (compiler.Engine, error) {
	logrus.Tracef("Creating %s compiler client from CLI configuration", constants.DriverKafka)
	return native.New(c)
}
