// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"time"

	"github.com/go-vela/pkg-queue/queue"
	"github.com/go-vela/server/random"
	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

// flags is a helper function to return the all
// supported command line interface (CLI) flags
// for the Worker.
func flags() []cli.Flag {
	// generate a new random string for the database encryption key
	//
	// https://pkg.go.dev/github.com/go-vela/server/random#GenerateRandomString
	key, err := random.GenerateRandomString(32)
	if err != nil {
		logrus.Fatal(err)
	}

	f := []cli.Flag{

		// Compiler Flags

		&cli.BoolFlag{
			EnvVars: []string{"VELA_COMPILER_GITHUB", "COMPILER_GITHUB"},
			Name:    "compiler.github.driver",
			Usage:   "enables the github compiler driver",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_COMPILER_GITHUB_ADDR", "COMPILER_GITHUB_ADDR"},
			Name:    "compiler.github.addr",
			Usage:   "fully qualified url (<scheme>://<host>) for compiler to pull github templates",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_COMPILER_GITHUB_TOKEN", "COMPILER_GITHUB_TOKEN"},
			Name:    "compiler.github.token",
			Usage:   "token used to access github system for compiler to pull github templates",
		},

		// Compiler Modification Flags

		&cli.StringFlag{
			EnvVars: []string{"VELA_COMPILER_MODIFICATION_ADDR", "COMPILER_MODIFICATION_ADDR"},
			Name:    "compiler.modification.addr",
			Usage:   "fully qualified url (<scheme>://<host>) for compiler to send requests to modification system",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_COMPILER_MODIFICATION_SECRET", "COMPILER_MODIFICATION_SECRET"},
			Name:    "compiler.modification.secret",
			Usage:   "secret used for communication between the compiler and modification system",
		},
		&cli.DurationFlag{
			EnvVars: []string{"VELA_COMPILER_MODIFICATION_DURATION", "COMPILER_MODIFICATION_DURATION"},
			Name:    "compiler.modification.duration",
			Usage:   "amount of time for requests sent to modification system",
			Value:   8 * time.Second,
		},
		&cli.IntFlag{
			EnvVars: []string{"VELA_COMPILER_MODIFICATION_RETRIES", "COMPILER_MODIFICATION_RETRIES"},
			Name:    "compiler.modification.retries",
			Usage:   "number of attempts compiler will resend requests to modification system",
			Value:   5,
		},

		// Database Flags

		&cli.StringFlag{
			EnvVars: []string{"VELA_DATABASE_DRIVER", "DATABASE_DRIVER"},
			Name:    "database.driver",
			Usage:   "driver to be used for the database",
			Value:   constants.DriverSqlite,
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_DATABASE_ADDR", "DATABASE_ADDR"},
			Name:    "database.addr",
			Usage:   "configuration string for the database",
			Value:   "vela.sqlite",
		},
		&cli.IntFlag{
			EnvVars: []string{"VELA_DATABASE_COMPRESSION_LEVEL", "DATABASE_COMPRESSION_LEVEL"},
			Name:    "database.compression.level",
			Usage:   "level of compression for logs stored in the database",
			Value:   constants.CompressionThree,
		},
		&cli.IntFlag{
			EnvVars: []string{"VELA_DATABASE_CONNECTION_IDLE", "DATABASE_CONNECTION_IDLE"},
			Name:    "database.connection.idle",
			Usage:   "number of idle connections to the database",
			Value:   2,
		},
		&cli.DurationFlag{
			EnvVars: []string{"VELA_DATABASE_CONNECTION_LIFE", "DATABASE_CONNECTION_LIFE"},
			Name:    "database.connection.life",
			Usage:   "amount of time a connection may be reused for the database",
			Value:   30 * time.Minute,
		},
		&cli.IntFlag{
			EnvVars: []string{"VELA_DATABASE_CONNECTION_OPEN", "DATABASE_CONNECTION_OPEN"},
			Name:    "database.connection.open",
			Usage:   "number of open connections to the database",
			Value:   0,
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_DATABASE_ENCRYPTION_KEY", "DATABASE_ENCRYPTION_KEY"},
			Name:    "database.encryption.key",
			Usage:   "32 character key for encrypting and decrypting data using AES-256",
			Value:   key,
		},

		// Logger Flags

		&cli.StringFlag{
			EnvVars: []string{"VELA_LOG_FORMAT", "LOG_FORMAT"},
			Name:    "log.format",
			Usage:   "log format to output",
			Value:   "json",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_LOG_LEVEL", "LOG_LEVEL"},
			Name:    "log.level",
			Usage:   "log level to output",
			Value:   "info",
		},

		// Metrics Flags

		&cli.DurationFlag{
			EnvVars: []string{"VELA_METRICS_WORKER_ACTIVE_DURATION", "METRICS_WORKER_ACTIVE_DURATION"},
			Name:    "metrics.worker.active.duration",
			Usage:   "amount of time for active workers within the /metrics endpoint",
			Value:   5 * time.Minute,
		},

		// Queue Flags

		&cli.StringFlag{
			EnvVars: []string{"VELA_QUEUE_DRIVER", "QUEUE_DRIVER"},
			Name:    "queue.driver",
			Usage:   "driver to be used for the queue",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_QUEUE_ADDR", "QUEUE_ADDR"},
			Name:    "queue.addr",
			Usage:   "fully qualified url (<scheme>://<host>) for the queue",
		},
		&cli.BoolFlag{
			EnvVars: []string{"VELA_QUEUE_CLUSTER", "QUEUE_CLUSTER"},
			Name:    "queue.cluster",
			Usage:   "enables connecting to a queue cluster",
		},
		&cli.StringSliceFlag{
			EnvVars: []string{"VELA_QUEUE_ROUTES", "QUEUE_ROUTES"},
			Name:    "queue.routes",
			Usage:   "list of routes (channels/topics) to publish builds",
			Value:   cli.NewStringSlice(constants.DefaultRoute),
		},

		// Secret Flags

		&cli.BoolFlag{
			EnvVars: []string{"VELA_SECRET_VAULT", "SECRET_VAULT"},
			Name:    "secret.vault.driver",
			Usage:   "enables the vault secret driver",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SECRET_VAULT_ADDR", "SECRET_VAULT_ADDR"},
			Name:    "secret.vault.addr",
			Usage:   "fully qualified url (<scheme>://<host>) for the vault system",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SECRET_VAULT_AUTH_METHOD", "SECRET_VAULT_AUTH_METHOD"},
			Name:    "secret.vault.auth-method",
			Usage:   "authentication method used to obtain token from vault system",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SECRET_VAULT_AWS_ROLE", "SECRET_VAULT_AWS_ROLE"},
			Name:    "secret.vault.aws-role",
			Usage:   "vault role used to connect to the auth/aws/login endpoint",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SECRET_VAULT_PREFIX", "SECRET_VAULT_PREFIX"},
			Name:    "secret.vault.prefix",
			Usage:   "prefix for k/v secrets in vault system e.g. secret/data/<prefix>/<path>",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SECRET_VAULT_TOKEN", "SECRET_VAULT_TOKEN"},
			Name:    "secret.vault.token",
			Usage:   "token used to access vault system",
		},
		&cli.DurationFlag{
			EnvVars: []string{"VELA_SECRET_VAULT_TOKEN_DURATION", "SECRET_VAULT_TOKEN_DURATION"},
			Name:    "secret.vault.token.duration",
			Usage:   "amount of time to wait before refreshing the vault token",
			Value:   30 * time.Minute,
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SECRET_VAULT_VERSION", "SECRET_VAULT_VERSION"},
			Name:    "secret.vault.version",
			Usage:   "version for the kv backend for the vault system",
			Value:   "2",
		},

		// Security Flags

		&cli.DurationFlag{
			EnvVars: []string{"VELA_ACCESS_TOKEN_DURATION", "ACCESS_TOKEN_DURATION"},
			Name:    "access.token.duration",
			Usage:   "amount of time an access token can be used",
			Value:   15 * time.Minute,
		},
		&cli.DurationFlag{
			EnvVars: []string{"VELA_REFRESH_TOKEN_DURATION", "REFRESH_TOKEN_DURATION"},
			Name:    "refresh.token.duration",
			Usage:   "amount of time a refresh token can be used",
			Value:   8 * time.Hour,
		},
		&cli.StringSliceFlag{
			EnvVars: []string{"VELA_REPO_ALLOW_LIST", "REPO_ALLOW_LIST"},
			Name:    "repo.allow.list",
			Usage:   "list of repos allowed to be activated in the system",
			Value:   &cli.StringSlice{},
		},
		&cli.BoolFlag{
			EnvVars: []string{"VELA_SECURE_COOKIE", "SECURE_COOKIE"},
			Name:    "secure.cookie",
			Usage:   "enables setting cookies with the secure flag - disable for testing",
			Value:   true,
		},
		&cli.BoolFlag{
			EnvVars: []string{"VELA_WEBHOOK_VALIDATION", "WEBHOOK_VALIDATION"},
			Name:    "webhook.validation",
			Usage:   "enables validating webhooks from the source system - disable for testing",
			Value:   true,
		},

		// Server Flags

		&cli.StringFlag{
			EnvVars: []string{"VELA_SERVER_ADDR", "SERVER_ADDR"},
			Name:    "server.addr",
			Usage:   "fully qualified url (<scheme>://<host>) for the server",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SERVER_PORT", "SERVER_PORT"},
			Name:    "server.port",
			Usage:   "port to publish server API and accept connections",
			Value:   "8080",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SERVER_SECRET", "SERVER_SECRET"},
			Name:    "server.secret",
			Usage:   "secret used for communication between the server and worker",
		},
		&cli.Int64Flag{
			EnvVars: []string{"VELA_SERVER_BUILD_TIMEOUT", "SERVER_BUILD_TIMEOUT"},
			Name:    "server.build.timeout",
			Usage:   "build timeout set by the server",
		},

		// Source Flags

		&cli.StringFlag{
			EnvVars: []string{"VELA_SOURCE_DRIVER", "SOURCE_DRIVER"},
			Name:    "source.driver",
			Usage:   "driver to be used for the version control system",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SOURCE_ADDR", "SOURCE_ADDR"},
			Name:    "source.addr",
			Usage:   "fully qualified url (<scheme>://<host>) for the version control system",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SOURCE_CLIENT", "SOURCE_CLIENT"},
			Name:    "source.client",
			Usage:   "OAuth client id from version control system",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SOURCE_SECRET", "SOURCE_SECRET"},
			Name:    "source.secret",
			Usage:   "OAuth client secret from version control system",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SOURCE_CONTEXT", "SOURCE_CONTEXT"},
			Name:    "source.context",
			Usage:   "context for commit status in version control system",
			Value:   "continuous-integration/vela",
		},

		// WebUI Flags

		&cli.StringFlag{
			EnvVars: []string{"VELA_WEBUI_ADDR", "WEBUI_ADDR"},
			Name:    "webui.addr",
			Usage:   "fully qualified url (<scheme>://<host>) for the web UI system",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_WEBUI_OAUTH_ENDPOINT", "WEBUI_OAUTH_ENDPOINT"},
			Name:    "webui.oauth.endpoint",
			Usage:   "endpoint path to web UI OAuth callback",
			Value:   "/account/authenticate",
		},
	}

	// Queue Flags

	f = append(f, queue.Flags...)

	return f
}
