// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/compiler/native"

	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
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
