// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/secret"
	"github.com/go-vela/server/source"
	"github.com/go-vela/server/version"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
)

// nolint: funlen // ignore function length due to flags
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

	// Database Flags

	app.Flags = append(app.Flags, database.Flags...)

	// Queue Flags

	app.Flags = append(app.Flags, queue.Flags...)

	// Secret Flags

	app.Flags = append(app.Flags, secret.Flags...)

	// Source Flags

	app.Flags = append(app.Flags, source.Flags...)

	// set logrus to log in JSON format
	logrus.SetFormatter(&logrus.JSONFormatter{})

	err = app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
