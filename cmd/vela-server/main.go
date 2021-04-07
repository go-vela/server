// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"os"
	"time"

	"github.com/go-vela/server/source"
	"github.com/go-vela/server/version"
	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// nolint: funlen // ignore function length due to flags
func main() {
	app := cli.NewApp()
	app.Name = "vela-server"
	app.Action = server
	app.Version = version.New().Semantic()

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			EnvVars: []string{"VELA_LOG_LEVEL", "LOG_LEVEL"},
			Name:    "log-level",
			Usage:   "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:   "info",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_ADDR", "VELA_HOST"},
			Name:    "server-addr",
			Usage:   "server address as a fully qualified url (<scheme>://<host>)",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_PORT"},
			Name:    "server-port",
			Usage:   "server port for the API to listen on",
			Value:   "8080",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_WEBUI_ADDR", "VELA_WEBUI_HOST"},
			Name:    "webui-addr",
			Usage:   "web ui address as a fully qualified url (<scheme>://<host>)",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_WEBUI_OAUTH_CALLBACK_PATH", "VELA_WEBUI_OAUTH_CALLBACK"},
			Name:    "webui-oauth-callback",
			Usage:   "web ui oauth callback path",
			Value:   "/account/authenticate",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SECRET"},
			Name:    "vela-secret",
			Usage:   "secret used for server <-> agent communication",
		},
		&cli.StringSliceFlag{
			EnvVars: []string{"VELA_REPO_ALLOWLIST"},
			Name:    "vela-repo-allowlist",
			Usage:   "allowlist is used to limit which repos can be activated within the system",
			Value:   &cli.StringSlice{},
		},
		&cli.BoolFlag{
			EnvVars: []string{"VELA_DISABLE_WEBHOOK_VALIDATION"},
			Name:    "vela-disable-webhook-validation",
			// nolint: lll // ignore long line length due to description
			Usage: "determines whether or not webhook validation is disabled.  useful for local development.",
			Value: false,
		},
		&cli.BoolFlag{
			EnvVars: []string{"VELA_ENABLE_SECURE_COOKIE"},
			Name:    "vela-enable-secure-cookie",
			Usage:   "determines whether or not use cookies with secure flag set.  useful for testing.",
			Value:   true,
		},
		&cli.Int64Flag{
			EnvVars: []string{"VELA_DEFAULT_BUILD_TIMEOUT"},
			Name:    "default-build-timeout",
			Usage:   "override default build timeout (minutes)",
		},

		// Security Flags

		&cli.DurationFlag{
			EnvVars: []string{"VELA_ACCESS_TOKEN_DURATION", "ACCESS_TOKEN_DURATION"},
			Name:    "access-token-duration",
			Usage:   "sets the duration of the access token",
			Value:   15 * time.Minute,
		},
		&cli.DurationFlag{
			EnvVars: []string{"VELA_REFRESH_TOKEN_DURATION", "REFRESH_TOKEN_DURATION"},
			Name:    "refresh-token-duration",
			Usage:   "sets the duration of the refresh token",
			Value:   8 * time.Hour,
		},

		// Compiler Flags

		&cli.BoolFlag{
			EnvVars: []string{"VELA_COMPILER_GITHUB", "COMPILER_GITHUB"},
			Name:    "github-driver",
			Usage:   "github compiler driver",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_COMPILER_GITHUB_URL", "COMPILER_GITHUB_URL"},
			Name:    "github-url",
			Usage:   "github url, used by compiler, for pulling registry templates",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_COMPILER_GITHUB_TOKEN", "COMPILER_GITHUB_TOKEN"},
			Name:    "github-token",
			Usage:   "github token, used by compiler, for pulling registry templates",
		},

		// Database Flags

		&cli.StringFlag{
			EnvVars: []string{"VELA_DATABASE_DRIVER", "DATABASE_DRIVER"},
			Name:    "database.driver",
			Usage:   "sets the driver to be used for the database",
			Value:   "sqlite3",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_DATABASE_CONFIG", "DATABASE_CONFIG"},
			Name:    "database.config",
			Usage:   "sets the configuration string to be used for the database",
			Value:   "vela.sqlite",
		},
		&cli.IntFlag{
			EnvVars: []string{"VELA_DATABASE_CONNECTION_OPEN", "DATABASE_CONNECTION_OPEN"},
			Name:    "database.connection.open",
			Usage:   "sets the number of open connections to the database",
			Value:   0,
		},
		&cli.IntFlag{
			EnvVars: []string{"VELA_DATABASE_CONNECTION_IDLE", "DATABASE_CONNECTION_IDLE"},
			Name:    "database.connection.idle",
			Usage:   "sets the number of idle connections to the database",
			Value:   2,
		},
		&cli.DurationFlag{
			EnvVars: []string{"VELA_DATABASE_CONNECTION_LIFE", "DATABASE_CONNECTION_LIFE"},
			Name:    "database.connection.life",
			Usage:   "sets the amount of time a connection may be reused for the database",
			Value:   30 * time.Minute,
		},
		&cli.IntFlag{
			EnvVars: []string{"VELA_DATABASE_COMPRESSION_LEVEL", "DATABASE_COMPRESSION_LEVEL"},
			Name:    "database.compression.level",
			Usage:   "sets the level of compression for logs stored in the database",
			Value:   constants.CompressionThree,
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_DATABASE_ENCRYPTION_KEY", "DATABASE_ENCRYPTION_KEY"},
			Name:    "database.encryption.key",
			Usage:   "AES-256 key for encrypting and decrypting values",
		},

		// Queue Flags

		&cli.StringFlag{
			EnvVars: []string{"VELA_QUEUE_DRIVER", "QUEUE_DRIVER"},
			Name:    "queue-driver",
			Usage:   "queue driver",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_QUEUE_CONFIG", "QUEUE_CONFIG"},
			Name:    "queue-config",
			Usage:   "queue driver configuration string",
		},
		&cli.BoolFlag{
			EnvVars: []string{"VELA_QUEUE_CLUSTER", "QUEUE_CLUSTER"},
			Name:    "queue-cluster",
			Usage:   "queue client is setup for clusters",
		},
		// By default all builds are pushed to the "vela" route
		&cli.StringSliceFlag{
			EnvVars: []string{"VELA_QUEUE_WORKER_ROUTES", "QUEUE_WORKER_ROUTES"},
			Name:    "queue-worker-routes",
			Usage:   "queue worker routes is configuration for routing builds",
		},

		// Secret Flags

		&cli.BoolFlag{
			EnvVars: []string{"VELA_SECRET_VAULT", "SECRET_VAULT"},
			Name:    "vault-driver",
			Usage:   "vault secret driver",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SECRET_VAULT_ADDR", "SECRET_VAULT_ADDR"},
			Name:    "vault-addr",
			Usage:   "vault address for storing secrets",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SECRET_VAULT_TOKEN", "SECRET_VAULT_TOKEN"},
			Name:    "vault-token",
			Usage:   "vault token for storing secrets",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SECRET_VAULT_VERSION", "SECRET_VAULT_VERSION"},
			Name:    "vault-version",
			Usage:   "vault k/v backend version to utilize",
			Value:   "2",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SECRET_VAULT_PREFIX", "SECRET_VAULT_PREFIX"},
			Name:    "vault-prefix",
			Usage:   "vault prefix for k/v secrets. e.g. secret/data/<prefix>/<path>",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SECRET_VAULT_AUTH_METHOD", "SECRET_VAULT_AUTH_METHOD"},
			Name:    "vault-auth-method",
			Usage:   "auth method to utilize to obtain token",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SECRET_VAULT_AWS_ROLE", "SECRET_VAULT_AWS_ROLE"},
			Name:    "vault-aws-role",
			Usage:   "vault role to connect to the auth/aws/login endpoint with",
		},
		&cli.DurationFlag{
			EnvVars: []string{"VELA_SECRET_VAULT_RENEWAL", "SECRET_VAULT_RENEWAL"},
			Name:    "vault-renewal",
			Usage:   "frequency which the vault token should be renewed",
			Value:   30 * time.Minute,
		},

		&cli.StringFlag{
			EnvVars: []string{"VELA_MODIFICATION_ADDR", "MODIFICATION_ADDR"},
			Name:    "modification-addr",
			Usage:   "modification address, used by compiler, endpoint to send pipeline for modification",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_MODIFICATION_SECRET", "MODIFICATION_SECRET"},
			Name:    "modification-secret",
			// nolint: lll // ignore long line length due to description
			Usage: "modification secret, used by compiler, secret to allow connectivity between compiler and modification endpoint",
		},
		&cli.DurationFlag{
			EnvVars: []string{"VELA_MODIFICATION_TIMEOUT", "MODIFICATION_TIMEOUT"},
			Name:    "modification-timeout",
			// nolint: lll // ignore long line length due to description
			Usage: "modification timeout, used by compiler, duration that the modification http request will timeout after",
			Value: 8 * time.Second,
		},
		&cli.IntFlag{
			EnvVars: []string{"VELA_MODIFICATION_RETRIES", "MODIFICATION_RETRIES"},
			Name:    "modification-retries",
			// nolint: lll // ignore long line length due to description
			Usage: "modification retries, used by compiler, number of http requires that the modification http request will fail after",
			Value: 5,
		},

		&cli.DurationFlag{
			EnvVars: []string{"VELA_WORKER_ACTIVE_INTERVAL", "WORKER_ACTIVE_INTERVAL"},
			Name:    "worker-active-interval",
			Usage:   "interval at which workers will show as active within the /metrics endpoint",
			Value:   5 * time.Minute,
		},
	}

	// Source Flags

	app.Flags = append(app.Flags, source.Flags...)

	// set logrus to log in JSON format
	logrus.SetFormatter(&logrus.JSONFormatter{})

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
