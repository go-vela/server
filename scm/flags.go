// SPDX-License-Identifier: Apache-2.0

package scm

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/util"
)

// Flags represents all supported command line
// interface (CLI) flags for the scm.
//
// https://pkg.go.dev/github.com/urfave/cli?tab=doc#Flag
var Flags = []cli.Flag{
	// SCM Flags

	&cli.StringFlag{
		Name:  "scm.driver",
		Usage: "driver to be used for the version control system",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SCM_DRIVER"),
			cli.EnvVar("SCM_DRIVER"),
			cli.File("/vela/scm/driver"),
		),
		Value: constants.DriverGithub,
	},
	&cli.StringFlag{
		Name:  "scm.addr",
		Usage: "fully qualified url (<scheme>://<host>) for the version control system",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SCM_ADDR"),
			cli.EnvVar("SCM_ADDR"),
			cli.File("/vela/scm/addr"),
		),
		Value: "https://github.com",
		Action: func(_ context.Context, _ *cli.Command, v string) error {
			if !strings.Contains(v, "://") {
				return fmt.Errorf("scm address must be fully qualified (<scheme>://<host>)")
			}

			if strings.HasSuffix(v, "/") {
				return fmt.Errorf("scm address must not have trailing slash")
			}

			return nil
		},
	},
	&cli.StringFlag{
		Name:  "scm.client",
		Usage: "OAuth client id from version control system",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SCM_CLIENT"),
			cli.EnvVar("SCM_CLIENT"),
			cli.File("/vela/scm/client"),
		),
		Required: true,
	},
	&cli.StringFlag{
		Name:  "scm.secret",
		Usage: "OAuth client secret from version control system",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SCM_SECRET"),
			cli.EnvVar("SCM_SECRET"),
			cli.File("/vela/scm/secret"),
		),
		Required: true,
	},
	&cli.BoolFlag{
		Name:    "vela-disable-webhook-validation",
		Usage:   "determines whether or not webhook validation is disabled.  useful for local development.",
		Sources: cli.EnvVars("VELA_DISABLE_WEBHOOK_VALIDATION"),
		Value:   false,
	},
	&cli.StringFlag{
		Name:  "scm.context",
		Usage: "context for commit status in version control system",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SCM_CONTEXT"),
			cli.EnvVar("SCM_CONTEXT"),
			cli.File("/vela/scm/context"),
		),
		Value: "continuous-integration/vela",
	},
	&cli.StringSliceFlag{
		Name:  "scm.scopes",
		Usage: "OAuth scopes to be used for the version control system",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SCM_SCOPES"),
			cli.EnvVar("SCM_SCOPES"),
			cli.File("/vela/scm/scopes"),
		),
		Value: []string{"repo", "repo:status", "user:email", "read:user", "read:org"},
	},
	&cli.StringFlag{
		Name: "scm.webhook.addr",
		Usage: "Alternative or proxy server address as a fully qualified url (<scheme>://<host>). " +
			"Use this when the Vela server address that the scm provider can send webhooks to " +
			"differs from the server address the UI and oauth flows use, such as when the server " +
			"is behind a Firewall or NAT, or when using something like ngrok to forward webhooks. " +
			"(defaults to VELA_ADDR).",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SCM_WEBHOOK_ADDR"),
			cli.EnvVar("SCM_WEBHOOK_ADDR"),
			cli.File("/vela/scm/webhook_addr"),
		),
	},
	&cli.Int64Flag{
		Name:  "scm.app.id",
		Usage: "set ID for the SCM App integration (GitHub App)",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SCM_APP_ID"),
			cli.EnvVar("SCM_APP_ID"),
			cli.File("/vela/scm/app_id"),
		),
		Action: func(_ context.Context, cmd *cli.Command, v int64) error {
			if v > 0 {
				if !cmd.Bool("vela-disable-webhook-validation") && cmd.String("scm.app.webhook-secret") == "" {
					return fmt.Errorf("webhook-validation enabled and app ID provided but no app webhook secret is provided")
				}

				if cmd.String("scm.app.private-key") == "" && cmd.String("scm.app.private-key.path") == "" {
					return fmt.Errorf("app ID provided but no app private key is provided")
				}

				if cmd.String("scm.app.private-key") != "" && cmd.String("scm.app.private-key.path") != "" {
					return fmt.Errorf("app ID provided but both app private key and app private key path are provided")
				}
			}

			return nil
		},
	},
	&cli.StringFlag{
		Name:  "scm.app.private-key",
		Usage: "set value of base64 encoded SCM App integration (GitHub App) private key",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SCM_APP_PRIVATE_KEY"),
			cli.EnvVar("SCM_APP_PRIVATE_KEY"),
			cli.File("/vela/scm/app_private_key"),
		),
	},
	&cli.StringFlag{
		Name:  "scm.app.private-key.path",
		Usage: "set filepath to the SCM App integration (GitHub App) private key",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SCM_APP_PRIVATE_KEY_PATH"),
			cli.EnvVar("SCM_APP_PRIVATE_KEY_PATH"),
			cli.File("/vela/scm/app_private_key_path"),
		),
	},
	&cli.StringFlag{
		Name:  "scm.app.webhook-secret",
		Usage: "set value of SCM App integration webhook secret",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SCM_APP_WEBHOOK_SECRET"),
			cli.EnvVar("SCM_APP_WEBHOOK_SECRET"),
			cli.File("/vela/scm/app_webhook_secret"),
		),
	},
	&cli.StringSliceFlag{
		Name:  "scm.app.permissions",
		Usage: "SCM App integration (GitHub App) permissions to be used as the allowed set of possible installation token permissions",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SCM_APP_PERMISSIONS"),
			cli.EnvVar("SCM_APP_PERMISSIONS"),
			cli.File("/vela/scm/app/permissions"),
		),
		Value: []string{"contents:read", "checks:write"},
	},
	&cli.StringMapFlag{
		Name:  "scm.repo.roles-map",
		Usage: "map of SCM roles to Vela permissions for repositories",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SCM_REPO_ROLES_MAP"),
			cli.EnvVar("SCM_REPO_ROLES_MAP"),
			cli.File("/vela/scm/repo/roles_map"),
		),
		Value: map[string]string{
			"admin":    constants.PermissionAdmin,
			"write":    constants.PermissionWrite,
			"maintain": constants.PermissionWrite,
			"triage":   constants.PermissionRead,
			"read":     constants.PermissionRead,
		},
		Action: func(_ context.Context, _ *cli.Command, v map[string]string) error {
			return util.ValidateRoleMap(v, "repo")
		},
	},
	&cli.StringMapFlag{
		Name:  "scm.org.roles-map",
		Usage: "map of SCM roles to Vela permissions for organizations",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SCM_ORG_ROLES_MAP"),
			cli.EnvVar("SCM_ORG_ROLES_MAP"),
			cli.File("/vela/scm/org/roles_map"),
		),
		Value: map[string]string{
			"admin":  constants.PermissionAdmin,
			"member": constants.PermissionRead,
		},
		Action: func(_ context.Context, _ *cli.Command, v map[string]string) error {
			return util.ValidateRoleMap(v, "org")
		},
	},
	&cli.StringMapFlag{
		Name:  "scm.team.roles-map",
		Usage: "map of SCM roles to Vela permissions for teams",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SCM_TEAM_ROLES_MAP"),
			cli.EnvVar("SCM_TEAM_ROLES_MAP"),
			cli.File("/vela/scm/team/roles_map"),
		),
		Value: map[string]string{
			"maintainer": constants.PermissionAdmin,
			"member":     constants.PermissionRead,
		},
		Action: func(_ context.Context, _ *cli.Command, v map[string]string) error {
			return util.ValidateRoleMap(v, "team")
		},
	},
}
