// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-vela/types/constants"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/secret"
	"github.com/go-vela/server/version"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
)

//nolint:funlen // ignore line length
func main() {
	// capture application version information
	v := version.New()

	// serialize the version information as pretty JSON
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logrus.Fatal(err)
	}

	// output the version information to stdout
	fmt.Fprintf(os.Stdout, "%s\n", string(bytes))

	app := cli.NewApp()
	app.Name = "vela-server"
	app.Action = server
	app.Version = v.Semantic()
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
		&cli.StringFlag{
			EnvVars: []string{"VELA_SERVER_PRIVATE_KEY"},
			Name:    "vela-server-private-key",
			Usage:   "private key used for signing tokens",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_CLONE_IMAGE"},
			Name:    "clone-image",
			Usage:   "the clone image to use for the injected clone step",
			Value:   "target/vela-git:v0.7.0@sha256:c2e8794556d6debceeaa2c82ff3cc9e8e6ed045b723419e3ff050409f25cc258",
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
			Usage:   "determines whether or not webhook validation is disabled.  useful for local development.",
			Value:   false,
		},
		&cli.BoolFlag{
			EnvVars: []string{"VELA_ENABLE_SECURE_COOKIE"},
			Name:    "vela-enable-secure-cookie",
			Usage:   "determines whether or not use cookies with secure flag set.  useful for testing.",
			Value:   true,
		},
		&cli.Int64Flag{
			EnvVars: []string{"VELA_DEFAULT_BUILD_LIMIT"},
			Name:    "default-build-limit",
			Usage:   "override default build limit",
			Value:   constants.BuildLimitDefault,
		},
		&cli.Int64Flag{
			EnvVars: []string{"VELA_MAX_BUILD_LIMIT"},
			Name:    "max-build-limit",
			Usage:   "override max build limit",
			Value:   constants.BuildLimitMax,
		},
		&cli.Int64Flag{
			EnvVars: []string{"VELA_DEFAULT_BUILD_TIMEOUT"},
			Name:    "default-build-timeout",
			Usage:   "override default build timeout (minutes)",
			Value:   constants.BuildTimeoutDefault,
		},
		&cli.StringSliceFlag{
			EnvVars: []string{"VELA_DEFAULT_REPO_EVENTS"},
			Name:    "default-repo-events",
			Usage:   "override default events for newly activated repositories",
			Value:   cli.NewStringSlice(constants.EventPush),
		},
		// Token Manager Flags
		&cli.DurationFlag{
			EnvVars: []string{"VELA_USER_ACCESS_TOKEN_DURATION", "USER_ACCESS_TOKEN_DURATION"},
			Name:    "user-access-token-duration",
			Usage:   "sets the duration of the user access token",
			Value:   15 * time.Minute,
		},
		&cli.DurationFlag{
			EnvVars: []string{"VELA_USER_REFRESH_TOKEN_DURATION", "USER_REFRESH_TOKEN_DURATION"},
			Name:    "user-refresh-token-duration",
			Usage:   "sets the duration of the user refresh token",
			Value:   8 * time.Hour,
		},
		&cli.DurationFlag{
			EnvVars: []string{"VELA_BUILD_TOKEN_BUFFER_DURATION", "BUILD_TOKEN_BUFFER_DURATION"},
			Name:    "build-token-buffer-duration",
			Usage:   "sets the duration of the buffer for build token expiration based on repo build timeout",
			Value:   5 * time.Minute,
		},
		&cli.DurationFlag{
			EnvVars: []string{"VELA_WORKER_AUTH_TOKEN_DURATION", "WORKER_AUTH_TOKEN_DURATION"},
			Name:    "worker-auth-token-duration",
			Usage:   "sets the duration of the worker auth token",
			Value:   20 * time.Minute,
		},
		&cli.DurationFlag{
			EnvVars: []string{"VELA_WORKER_REGISTER_TOKEN_DURATION", "WORKER_REGISTER_TOKEN_DURATION"},
			Name:    "worker-register-token-duration",
			Usage:   "sets the duration of the worker register token",
			Value:   1 * time.Minute,
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
		&cli.StringFlag{
			EnvVars: []string{"VELA_MODIFICATION_ADDR", "MODIFICATION_ADDR"},
			Name:    "modification-addr",
			Usage:   "modification address, used by compiler, endpoint to send pipeline for modification",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_MODIFICATION_SECRET", "MODIFICATION_SECRET"},
			Name:    "modification-secret",
			Usage:   "modification secret, used by compiler, secret to allow connectivity between compiler and modification endpoint",
		},
		&cli.DurationFlag{
			EnvVars: []string{"VELA_MODIFICATION_TIMEOUT", "MODIFICATION_TIMEOUT"},
			Name:    "modification-timeout",
			Usage:   "modification timeout, used by compiler, duration that the modification http request will timeout after",
			Value:   8 * time.Second,
		},
		&cli.IntFlag{
			EnvVars: []string{"VELA_MODIFICATION_RETRIES", "MODIFICATION_RETRIES"},
			Name:    "modification-retries",
			Usage:   "modification retries, used by compiler, number of http requires that the modification http request will fail after",
			Value:   5,
		},
		&cli.DurationFlag{
			EnvVars: []string{"VELA_WORKER_ACTIVE_INTERVAL", "WORKER_ACTIVE_INTERVAL"},
			Name:    "worker-active-interval",
			Usage:   "interval at which workers will show as active within the /metrics endpoint",
			Value:   5 * time.Minute,
		},
	}
	// Add Database Flags
	app.Flags = append(app.Flags, database.Flags...)

	// Add Queue Flags
	app.Flags = append(app.Flags, queue.Flags...)

	// Add Secret Flags
	app.Flags = append(app.Flags, secret.Flags...)

	// Add Source Flags
	app.Flags = append(app.Flags, scm.Flags...)

	// set logrus to log in JSON format
	logrus.SetFormatter(&logrus.JSONFormatter{})

	if err = app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
