// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/secret"
	"github.com/go-vela/server/tracing"
	"github.com/go-vela/server/version"
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
			EnvVars: []string{"VELA_LOG_FORMATTER", "LOG_FORMATTER"},
			Name:    "log-formatter",
			Usage:   "set log formatter - options: (json|ecs)",
			Value:   "json",
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
		&cli.StringSliceFlag{
			EnvVars: []string{"VELA_CORS_ALLOW_ORIGINS", "VELA_CORS_ALLOWED_ORIGINS"},
			Name:    "cors-allow-origins",
			Usage:   "list of origins a cross-domain request can be executed from",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_SECRET"},
			Name:    "vela-secret",
			Usage:   "secret used for server <-> agent communication",
		},
		&cli.DurationFlag{
			EnvVars: []string{"VELA_PLATFORM_SETTINGS_REFRESH_INTERVAL", "VELA_SETTINGS_REFRESH_INTERVAL"},
			Name:    "settings-refresh-interval",
			Usage:   "interval at which platform settings will be refreshed",
			Value:   5 * time.Second,
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
			Value:   "target/vela-git-slim:v0.12.1@sha256:93cdb399e0a3150addac494198473c464c978ca055121593607097b75480192b", // renovate: container
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
		&cli.Int64Flag{
			EnvVars: []string{"VELA_DEFAULT_APPROVAL_TIMEOUT"},
			Name:    "default-approval-timeout",
			Usage:   "override default approval timeout (days)",
			Value:   constants.ApprovalTimeoutDefault,
		},
		&cli.StringSliceFlag{
			EnvVars: []string{"VELA_DEFAULT_REPO_EVENTS"},
			Name:    "default-repo-events",
			Usage:   "override default events for newly activated repositories",
			Value:   cli.NewStringSlice(constants.EventPush),
		},
		&cli.Int64Flag{
			EnvVars: []string{"VELA_DEFAULT_REPO_EVENTS_MASK"},
			Name:    "default-repo-events-mask",
			Usage:   "set default event mask for newly activated repositories",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_DEFAULT_REPO_APPROVE_BUILD"},
			Name:    "default-repo-approve-build",
			Usage:   "override default approve build for newly activated repositories",
			Value:   constants.ApproveForkAlways,
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
		&cli.DurationFlag{
			EnvVars: []string{"VELA_OPEN_ID_TOKEN_DURATION", "OPEN_ID_TOKEN_DURATION"},
			Name:    "id-token-duration",
			Usage:   "sets the duration of an OpenID token requested during a build (should be short)",
			Value:   5 * time.Minute,
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_OPEN_ID_ISSUER", "OPEN_ID_ISSUER"},
			Name:    "oidc-issuer",
			Usage:   "sets the issuer of the OpenID token requested during a build",
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
		&cli.Uint64Flag{
			EnvVars: []string{"VELA_COMPILER_STARLARK_EXEC_LIMIT", "COMPILER_STARLARK_EXEC_LIMIT"},
			Name:    "compiler-starlark-exec-limit",
			Usage:   "set the starlark execution step limit for compiling starlark pipelines",
			Value:   7500,
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
		&cli.IntFlag{
			EnvVars: []string{"VELA_MAX_TEMPLATE_DEPTH", "MAX_TEMPLATE_DEPTH"},
			Name:    "max-template-depth",
			Usage:   "max template depth, used by compiler, maximum number of templates that can be called in a template chain",
			Value:   3,
		},
		&cli.DurationFlag{
			EnvVars: []string{"VELA_WORKER_ACTIVE_INTERVAL", "WORKER_ACTIVE_INTERVAL"},
			Name:    "worker-active-interval",
			Usage:   "interval at which workers will show as active within the /metrics endpoint",
			Value:   5 * time.Minute,
		},
		// schedule flags
		&cli.DurationFlag{
			EnvVars: []string{"VELA_SCHEDULE_MINIMUM_FREQUENCY", "SCHEDULE_MINIMUM_FREQUENCY"},
			Name:    "schedule-minimum-frequency",
			Usage:   "minimum time allowed between each build triggered for a schedule",
			Value:   1 * time.Hour,
		},
		&cli.DurationFlag{
			EnvVars: []string{"VELA_SCHEDULE_INTERVAL", "SCHEDULE_INTERVAL"},
			Name:    "schedule-interval",
			Usage:   "interval at which schedules will be processed by the server to trigger builds",
			Value:   5 * time.Minute,
		},
		&cli.StringSliceFlag{
			EnvVars: []string{"VELA_SCHEDULE_ALLOWLIST"},
			Name:    "vela-schedule-allowlist",
			Usage:   "limit which repos can be utilize the schedule feature within the system",
			Value:   &cli.StringSlice{},
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

	// Add Tracing Flags
	app.Flags = append(app.Flags, tracing.Flags...)

	if err = app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
