// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"os"
	"time"

	"github.com/go-vela/server/version"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "vela-server"
	app.Action = server
	app.Version = version.Version.String()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			EnvVar: "VELA_PORT",
			Name:   "server-port",
			Usage:  "API port to listen on",
			Value:  ":8080",
		},
		cli.StringFlag{
			EnvVar: "VELA_LOG_LEVEL,LOG_LEVEL",
			Name:   "log-level",
			Usage:  "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:  "info",
		},
		cli.StringFlag{
			EnvVar: "VELA_ADDR,VELA_HOST",
			Name:   "server-addr",
			Usage:  "server address as a fully qualified url (<scheme>://<host>)",
		},
		cli.StringFlag{
			EnvVar: "VELA_SECRET",
			Name:   "vela-secret",
			Usage:  "secret used for server <-> agent communication",
		},
		cli.StringSliceFlag{
			EnvVar: "VELA_REPO_WHITELIST",
			Name:   "vela-repo-whitelist",
			Usage:  "whitelist is used to limit which repos can be activated within the system",
			Value:  &cli.StringSlice{},
		},

		// Compiler Flags

		cli.BoolFlag{
			EnvVar: "VELA_COMPILER_GITHUB,COMPILER_GITHUB",
			Name:   "github-driver",
			Usage:  "github compiler driver",
		},
		cli.StringFlag{
			EnvVar: "VELA_COMPILER_GITHUB_URL,COMPILER_GITHUB_URL",
			Name:   "github-url",
			Usage:  "github url, used by compiler, for pulling registry templates",
		},
		cli.StringFlag{
			EnvVar: "VELA_COMPILER_GITHUB_TOKEN,COMPILER_GITHUB_TOKEN",
			Name:   "github-token",
			Usage:  "github token, used by compiler, for pulling registry templates",
		},

		// Database Flags

		cli.StringFlag{
			EnvVar: "VELA_DATABASE_DRIVER,DATABASE_DRIVER",
			Name:   "database.driver",
			Usage:  "sets the driver to be used for the database",
			Value:  "sqlite3",
		},
		cli.StringFlag{
			EnvVar: "VELA_DATABASE_CONFIG,DATABASE_CONFIG",
			Name:   "database.config",
			Usage:  "sets the configuration string to be used for the database",
			Value:  "vela.sqlite",
		},
		cli.IntFlag{
			EnvVar: "VELA_DATABASE_CONNECTION_OPEN,DATABASE_CONNECTION_OPEN",
			Name:   "database.connection.open",
			Usage:  "sets the number of open connections to the database",
			Value:  0,
		},
		cli.IntFlag{
			EnvVar: "VELA_DATABASE_CONNECTION_IDLE,DATABASE_CONNECTION_IDLE",
			Name:   "database.connection.idle",
			Usage:  "sets the number of idle connections to the database",
			Value:  2,
		},
		cli.DurationFlag{
			EnvVar: "VELA_DATABASE_CONNECTION_LIFE,DATABASE_CONNECTION_LIFE",
			Name:   "database.connection.life",
			Usage:  "sets the amount of time a connection may be reused for the database",
			Value:  30 * time.Minute,
		},

		// Queue Flags

		cli.StringFlag{
			EnvVar: "VELA_QUEUE_DRIVER,QUEUE_DRIVER",
			Name:   "queue-driver",
			Usage:  "queue driver",
		},
		cli.StringFlag{
			EnvVar: "VELA_QUEUE_CONFIG,QUEUE_CONFIG",
			Name:   "queue-config",
			Usage:  "queue driver configuration string",
		},
		cli.BoolFlag{
			EnvVar: "VELA_QUEUE_CLUSTER,QUEUE_CLUSTER",
			Name:   "queue-cluster",
			Usage:  "queue client is setup for clusters",
		},

		// Secret Flags

		cli.BoolFlag{
			EnvVar: "VELA_SECRET_VAULT,SECRET_VAULT",
			Name:   "vault-driver",
			Usage:  "vault secret driver",
		},
		cli.StringFlag{
			EnvVar: "VELA_SECRET_VAULT_ADDR,SECRET_VAULT_ADDR",
			Name:   "vault-addr",
			Usage:  "vault address for storing secrets",
		},
		cli.StringFlag{
			EnvVar: "VELA_SECRET_VAULT_TOKEN,SECRET_VAULT_TOKEN",
			Name:   "vault-token",
			Usage:  "vault token for storing secrets",
		},

		// Source Flags

		cli.StringFlag{
			EnvVar: "VELA_SOURCE_DRIVER,SOURCE_DRIVER",
			Name:   "source-driver",
			Usage:  "source driver",
		},
		cli.StringFlag{
			EnvVar: "VELA_SOURCE_URL,SOURCE_URL",
			Name:   "source-url",
			Usage:  "source url address",
		},
		cli.StringFlag{
			EnvVar: "VELA_SOURCE_CLIENT,SOURCE_CLIENT",
			Name:   "source-client",
			Usage:  "source client id",
		},
		cli.StringFlag{
			EnvVar: "VELA_SOURCE_SECRET,SOURCE_SECRET",
			Name:   "source-secret",
			Usage:  "source client secret",
		},
		cli.StringFlag{
			EnvVar: "VELA_SOURCE_CONTEXT,SOURCE_CONTEXT",
			Name:   "source-context",
			Usage:  "source commit status context",
			Value:  "continuous-integration/vela",
		},
	}

	// set logrus to log in JSON format
	logrus.SetFormatter(&logrus.JSONFormatter{})

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
