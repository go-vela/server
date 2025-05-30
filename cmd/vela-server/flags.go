// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/internal/image"
)

// Flags represents all supported command line
// interface (CLI) flags for the core server.
//
// https://pkg.go.dev/github.com/urfave/cli?tab=doc#Flag
var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "log-level",
		Usage:   "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
		Sources: cli.EnvVars("VELA_LOG_LEVEL", "LOG_LEVEL"),
		Value:   "info",
	},
	&cli.StringFlag{
		Name:    "log-formatter",
		Usage:   "set log formatter - options: (json|ecs)",
		Sources: cli.EnvVars("VELA_LOG_FORMATTER", "LOG_FORMATTER"),
		Value:   "json",
	},
	&cli.StringFlag{
		Name:     "server-addr",
		Usage:    "server address as a fully qualified url (<scheme>://<host>)",
		Required: true,
		Sources:  cli.EnvVars("VELA_ADDR", "VELA_HOST"),
		Action: func(_ context.Context, cmd *cli.Command, v string) error {
			if !strings.Contains(v, "://") {
				return fmt.Errorf("server-addr (VELA_ADDR or VELA_HOST) flag must be <scheme>://<hostname> format")
			}

			if strings.HasSuffix(v, "/") {
				return fmt.Errorf("server-addr (VELA_ADDR or VELA_HOST) flag must not have trailing slash")
			}

			// warn if corresponding webui addr is not set
			if len(cmd.String("webui-addr")) == 0 {
				logrus.Warn("optional flag webui-addr (VELA_WEBUI_ADDR or VELA_WEBUI_HOST) not set")

				return nil
			}

			return nil
		},
	},
	&cli.StringFlag{
		Name:    "server-port",
		Usage:   "server port for the API to listen on",
		Sources: cli.EnvVars("VELA_PORT"),
		Value:   "8080",
	},
	&cli.StringFlag{
		Name:    "webui-addr",
		Usage:   "web ui address as a fully qualified url (<scheme>://<host>)",
		Sources: cli.EnvVars("VELA_WEBUI_ADDR", "VELA_WEBUI_HOST"),
		Action: func(_ context.Context, _ *cli.Command, v string) error {
			if !strings.Contains(v, "://") {
				return fmt.Errorf("webui-addr (VELA_WEBUI_ADDR or VELA_WEBUI_HOST) flag must be <scheme>://<hostname> format")
			}

			if strings.HasSuffix(v, "/") {
				return fmt.Errorf("webui-addr (VELA_WEBUI_ADDR or VELA_WEBUI_HOST) flag must not have trailing slash")
			}

			return nil
		},
	},
	&cli.StringFlag{
		Name:    "webui-oauth-callback",
		Usage:   "web ui oauth callback path",
		Sources: cli.EnvVars("VELA_WEBUI_OAUTH_CALLBACK_PATH", "VELA_WEBUI_OAUTH_CALLBACK"),
		Value:   "/account/authenticate",
	},
	&cli.StringSliceFlag{
		Name:    "cors-allow-origins",
		Usage:   "list of origins a cross-domain request can be executed from",
		Sources: cli.EnvVars("VELA_CORS_ALLOW_ORIGINS", "CORS_ALLOW_ORIGINS"),
	},
	&cli.StringFlag{
		Name:    "vela-secret",
		Usage:   "secret used for server <-> agent communication",
		Sources: cli.EnvVars("VELA_SECRET"),
	},
	&cli.DurationFlag{
		Name:    "settings-refresh-interval",
		Usage:   "interval at which platform settings will be refreshed",
		Sources: cli.EnvVars("VELA_PLATFORM_SETTINGS_REFRESH_INTERVAL", "VELA_SETTINGS_REFRESH_INTERVAL"),
		Value:   5 * time.Second,
	},
	&cli.StringFlag{
		Name:     "vela-server-private-key",
		Usage:    "private key used for signing tokens",
		Required: true,
		Sources:  cli.EnvVars("VELA_SERVER_PRIVATE_KEY"),
	},
	&cli.StringFlag{
		Name:    "clone-image",
		Usage:   "the clone image to use for the injected clone step",
		Sources: cli.EnvVars("VELA_CLONE_IMAGE"),
		Value:   "target/vela-git-slim:v0.12.1@sha256:93cdb399e0a3150addac494198473c464c978ca055121593607097b75480192b", // renovate: container
		Action: func(_ context.Context, _ *cli.Command, v string) error {
			_, err := image.ParseWithError(v)
			if err != nil {
				return fmt.Errorf("invalid clone image %s: %w", v, err)
			}

			return nil
		},
	},
	&cli.StringSliceFlag{
		Name:    "vela-repo-allowlist",
		Usage:   "allowlist is used to limit which repos can be activated within the system",
		Sources: cli.EnvVars("VELA_REPO_ALLOWLIST"),
		Value:   []string{},
	},
	&cli.BoolFlag{
		Name:    "vela-enable-secure-cookie",
		Usage:   "determines whether or not use cookies with secure flag set.  useful for testing.",
		Sources: cli.EnvVars("VELA_ENABLE_SECURE_COOKIE"),
		Value:   true,
	},
	&cli.Int32Flag{
		Name:    "default-build-limit",
		Usage:   "override default build limit",
		Sources: cli.EnvVars("VELA_DEFAULT_BUILD_LIMIT"),
		Value:   constants.BuildLimitDefault,
		Action: func(_ context.Context, _ *cli.Command, v int32) error {
			if v <= 0 {
				return fmt.Errorf("default-build-limit (VELA_DEFAULT_BUILD_LIMIT) flag must be greater than 0")
			}

			return nil
		},
	},
	&cli.Int32Flag{
		Name:    "max-build-limit",
		Usage:   "override max build limit",
		Sources: cli.EnvVars("VELA_MAX_BUILD_LIMIT"),
		Value:   constants.BuildLimitMax,
		Action: func(_ context.Context, cmd *cli.Command, v int32) error {
			if v <= 0 {
				return fmt.Errorf("max-build-limit (VELA_MAX_BUILD_LIMIT) flag must be greater than 0")
			}

			if cmd.Int32("default-build-limit") > v {
				return fmt.Errorf("max-build-limit (VELA_MAX_BUILD_LIMIT) must be greater than default-build-limit (VELA_DEFAULT_BUILD_LIMIT)")
			}

			return nil
		},
	},
	&cli.Int32Flag{
		Name:    "default-build-timeout",
		Usage:   "override default build timeout (minutes)",
		Sources: cli.EnvVars("VELA_DEFAULT_BUILD_TIMEOUT"),
		Value:   constants.BuildTimeoutDefault,
	},
	&cli.Int32Flag{
		Name:    "default-approval-timeout",
		Usage:   "override default approval timeout (days)",
		Sources: cli.EnvVars("VELA_DEFAULT_APPROVAL_TIMEOUT"),
		Value:   constants.ApprovalTimeoutDefault,
	},
	&cli.StringSliceFlag{
		Name:    "default-repo-events",
		Usage:   "override default events for newly activated repositories",
		Sources: cli.EnvVars("VELA_DEFAULT_REPO_EVENTS"),
		Value:   []string{constants.EventPush},
		Action: func(_ context.Context, _ *cli.Command, v []string) error {
			for _, event := range v {
				switch event {
				case constants.EventPull:
				case constants.EventPush:
				case constants.EventDeploy:
				case constants.EventTag:
				case constants.EventComment:
				default:
					return fmt.Errorf("default-repo-events (VELA_DEFAULT_REPO_EVENTS) has the unsupported value of %s", event)
				}
			}

			return nil
		},
	},
	&cli.Int64Flag{
		Name:    "default-repo-events-mask",
		Usage:   "set default event mask for newly activated repositories",
		Sources: cli.EnvVars("VELA_DEFAULT_REPO_EVENTS_MASK"),
	},
	&cli.StringFlag{
		Name:    "default-repo-approve-build",
		Usage:   "override default approve build for newly activated repositories",
		Sources: cli.EnvVars("VELA_DEFAULT_REPO_APPROVE_BUILD"),
		Value:   constants.ApproveForkAlways,
		Action: func(_ context.Context, _ *cli.Command, v string) error {
			if v != constants.ApproveForkAlways &&
				v != constants.ApproveNever &&
				v != constants.ApproveForkNoWrite &&
				v != constants.ApproveOnce {
				return fmt.Errorf("default-repo-approve-build (VELA_DEFAULT_REPO_APPROVE_BUILD) has the unsupported value of %s", v)
			}

			return nil
		},
	},
	&cli.Int32Flag{
		Name:    "max-dashboard-repos",
		Usage:   "set the maximum amount of repos that can belong to a dashboard",
		Sources: cli.EnvVars("VELA_MAX_DASHBOARD_REPOS"),
		Value:   10,
	},
	&cli.Int32Flag{
		Name:    "queue-restart-limit",
		Usage:   "set the max queue size before pending builds are no longer allowed to be restarted (set to 0 to not enforce)",
		Sources: cli.EnvVars("VELA_QUEUE_RESTART_LIMIT"),
		Value:   30,
	},
	// Token Manager Flags
	&cli.DurationFlag{
		Name:    "user-access-token-duration",
		Usage:   "sets the duration of the user access token",
		Sources: cli.EnvVars("VELA_USER_ACCESS_TOKEN_DURATION", "USER_ACCESS_TOKEN_DURATION"),
		Value:   15 * time.Minute,
	},
	&cli.DurationFlag{
		Name:    "user-refresh-token-duration",
		Usage:   "sets the duration of the user refresh token",
		Sources: cli.EnvVars("VELA_USER_REFRESH_TOKEN_DURATION", "USER_REFRESH_TOKEN_DURATION"),
		Value:   8 * time.Hour,
		Action: func(_ context.Context, cmd *cli.Command, v time.Duration) error {
			if cmd.Duration("user-access-token-duration").Seconds() >= v.Seconds() {
				return fmt.Errorf("user-refresh-token-duration (VELA_USER_REFRESH_TOKEN_DURATION) must be larger than the user-access-token-duration (VELA_USER_ACCESS_TOKEN_DURATION)")
			}

			return nil
		},
	},
	&cli.DurationFlag{
		Name:    "build-token-buffer-duration",
		Usage:   "sets the duration of the buffer for build token expiration based on repo build timeout",
		Sources: cli.EnvVars("VELA_BUILD_TOKEN_BUFFER_DURATION", "BUILD_TOKEN_BUFFER_DURATION"),
		Value:   5 * time.Minute,
		Action: func(_ context.Context, _ *cli.Command, v time.Duration) error {
			if v.Seconds() <= 0 {
				return fmt.Errorf("build-token-buffer-duration (VELA_BUILD_TOKEN_BUFFER_DURATION) must not be a negative time value")
			}

			return nil
		},
	},
	&cli.DurationFlag{
		Name:    "worker-auth-token-duration",
		Usage:   "sets the duration of the worker auth token",
		Sources: cli.EnvVars("VELA_WORKER_AUTH_TOKEN_DURATION", "WORKER_AUTH_TOKEN_DURATION"),
		Value:   20 * time.Minute,
	},
	&cli.DurationFlag{
		Name:    "worker-register-token-duration",
		Usage:   "sets the duration of the worker register token",
		Sources: cli.EnvVars("VELA_WORKER_REGISTER_TOKEN_DURATION", "WORKER_REGISTER_TOKEN_DURATION"),
		Value:   1 * time.Minute,
	},
	&cli.DurationFlag{
		Name:    "id-token-duration",
		Usage:   "sets the duration of an OpenID token requested during a build (should be short)",
		Sources: cli.EnvVars("VELA_OPEN_ID_TOKEN_DURATION", "OPEN_ID_TOKEN_DURATION"),
		Value:   5 * time.Minute,
	},
	&cli.StringFlag{
		Name:    "oidc-issuer",
		Usage:   "sets the issuer of the OpenID token requested during a build",
		Sources: cli.EnvVars("VELA_OPEN_ID_ISSUER", "OPEN_ID_ISSUER"),
		Action: func(_ context.Context, _ *cli.Command, v string) error {
			if len(v) > 0 {
				_, err := url.Parse(v)
				if err != nil {
					return fmt.Errorf("oidc-issuer (VELA_OPEN_ID_ISSUER) flag must be a valid URL")
				}
			}

			return nil
		},
	},
	// Compiler Flags
	&cli.BoolFlag{
		Name:    "github-driver",
		Usage:   "github compiler driver",
		Sources: cli.EnvVars("VELA_COMPILER_GITHUB", "COMPILER_GITHUB"),
		Action: func(_ context.Context, cmd *cli.Command, v bool) error {
			if v && len(cmd.String("github-url")) == 0 {
				return fmt.Errorf("github-url (VELA_COMPILER_GITHUB_URL or COMPILER_GITHUB_URL) flag not specified")
			}

			if v && len(cmd.String("github-token")) == 0 {
				return fmt.Errorf("github-token (VELA_COMPILER_GITHUB_TOKEN or COMPILER_GITHUB_TOKEN) flag not specified")
			}

			return nil
		},
	},
	&cli.StringFlag{
		Name:    "github-url",
		Usage:   "github url, used by compiler, for pulling registry templates",
		Sources: cli.EnvVars("VELA_COMPILER_GITHUB_URL", "COMPILER_GITHUB_URL"),
	},
	&cli.StringFlag{
		Name:    "github-token",
		Usage:   "github token, used by compiler, for pulling registry templates",
		Sources: cli.EnvVars("VELA_COMPILER_GITHUB_TOKEN", "COMPILER_GITHUB_TOKEN"),
	},
	&cli.Int64Flag{
		Name:    "compiler-starlark-exec-limit",
		Usage:   "set the starlark execution step limit for compiling starlark pipelines",
		Sources: cli.EnvVars("VELA_COMPILER_STARLARK_EXEC_LIMIT", "COMPILER_STARLARK_EXEC_LIMIT"),
		Value:   7500,
	},
	&cli.StringFlag{
		Name:    "modification-addr",
		Usage:   "modification address, used by compiler, endpoint to send pipeline for modification",
		Sources: cli.EnvVars("VELA_MODIFICATION_ADDR", "MODIFICATION_ADDR"),
	},
	&cli.StringFlag{
		Name:    "modification-secret",
		Usage:   "modification secret, used by compiler, secret to allow connectivity between compiler and modification endpoint",
		Sources: cli.EnvVars("VELA_MODIFICATION_SECRET", "MODIFICATION_SECRET"),
	},
	&cli.DurationFlag{
		Name:    "modification-timeout",
		Usage:   "modification timeout, used by compiler, duration that the modification http request will timeout after",
		Sources: cli.EnvVars("VELA_MODIFICATION_TIMEOUT", "MODIFICATION_TIMEOUT"),
		Value:   8 * time.Second,
	},
	&cli.IntFlag{
		Name:    "modification-retries",
		Usage:   "modification retries, used by compiler, number of http requires that the modification http request will fail after",
		Sources: cli.EnvVars("VELA_MODIFICATION_RETRIES", "MODIFICATION_RETRIES"),
		Value:   5,
	},
	&cli.IntFlag{
		Name:    "max-template-depth",
		Usage:   "max template depth, used by compiler, maximum number of templates that can be called in a template chain",
		Sources: cli.EnvVars("VELA_MAX_TEMPLATE_DEPTH", "MAX_TEMPLATE_DEPTH"),
		Value:   3,
		Action: func(_ context.Context, _ *cli.Command, v int) error {
			if v < 1 {
				return fmt.Errorf("max-template-depth (VELA_MAX_TEMPLATE_DEPTH) or (MAX_TEMPLATE_DEPTH) flag must be greater than 0")
			}

			return nil
		},
	},
	&cli.DurationFlag{
		Name:    "worker-active-interval",
		Usage:   "interval at which workers will show as active within the /metrics endpoint",
		Sources: cli.EnvVars("VELA_WORKER_ACTIVE_INTERVAL", "WORKER_ACTIVE_INTERVAL"),
		Value:   5 * time.Minute,
	},
	// schedule flags
	&cli.DurationFlag{
		Name:    "schedule-minimum-frequency",
		Usage:   "minimum time allowed between each build triggered for a schedule",
		Sources: cli.EnvVars("VELA_SCHEDULE_MINIMUM_FREQUENCY", "SCHEDULE_MINIMUM_FREQUENCY"),
		Value:   1 * time.Hour,
	},
	&cli.DurationFlag{
		Name:    "schedule-interval",
		Usage:   "interval at which schedules will be processed by the server to trigger builds",
		Sources: cli.EnvVars("VELA_SCHEDULE_INTERVAL", "SCHEDULE_INTERVAL"),
		Value:   5 * time.Minute,
	},
	&cli.StringSliceFlag{
		Name:    "vela-schedule-allowlist",
		Usage:   "limit which repos can be utilize the schedule feature within the system",
		Sources: cli.EnvVars("VELA_SCHEDULE_ALLOWLIST"),
		Value:   []string{},
	},
}
