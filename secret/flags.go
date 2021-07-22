// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"time"

	"github.com/urfave/cli/v2"
)

// Flags represents all supported command line
// interface (CLI) flags for the secret.
//
// https://pkg.go.dev/github.com/urfave/cli?tab=doc#Flag
var Flags = []cli.Flag{

	// Logger Flags

	&cli.StringFlag{
		EnvVars:  []string{"VELA_SECRET_LOG_FORMAT", "SECRET_LOG_FORMAT", "VELA_LOG_FORMAT"},
		FilePath: "/vela/secret/log_format",
		Name:     "secret.log.format",
		Usage:    "format of logs to output",
		Value:    "json",
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_SECRET_LOG_LEVEL", "SECRET_LOG_LEVEL", "VELA_LOG_LEVEL"},
		FilePath: "/vela/secret/log_level",
		Name:     "secret.log.level",
		Usage:    "level of logs to output",
		Value:    "info",
	},

	// Secret Flags

	&cli.BoolFlag{
		EnvVars:  []string{"VELA_SECRET_VAULT", "SECRET_VAULT"},
		FilePath: "/vela/secret/vault/driver",
		Name:     "secret.vault.driver",
		Usage:    "enables the vault secret driver",
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_SECRET_VAULT_ADDR", "SECRET_VAULT_ADDR"},
		FilePath: "/vela/secret/vault/addr",
		Name:     "secret.vault.addr",
		Usage:    "fully qualified url (<scheme>://<host>) for the vault system",
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_SECRET_VAULT_AUTH_METHOD", "SECRET_VAULT_AUTH_METHOD"},
		FilePath: "/vela/secret/vault/auth_method",
		Name:     "secret.vault.auth-method",
		Usage:    "authentication method used to obtain token from vault system",
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_SECRET_VAULT_AWS_ROLE", "SECRET_VAULT_AWS_ROLE"},
		FilePath: "/vela/secret/vault/aws_role",
		Name:     "secret.vault.aws-role",
		Usage:    "vault role used to connect to the auth/aws/login endpoint",
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_SECRET_VAULT_PREFIX", "SECRET_VAULT_PREFIX"},
		FilePath: "/vela/secret/vault/prefix",
		Name:     "secret.vault.prefix",
		Usage:    "prefix for k/v secrets in vault system e.g. secret/data/<prefix>/<path>",
	},
	&cli.DurationFlag{
		EnvVars:  []string{"VELA_SECRET_VAULT_RENEWAL", "SECRET_VAULT_RENEWAL"},
		FilePath: "/vela/secret/vault/renewal",
		Name:     "secret.vault.renewal",
		Usage:    "frequency which the vault token should be renewed",
		Value:    30 * time.Minute,
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_SECRET_VAULT_TOKEN", "SECRET_VAULT_TOKEN"},
		FilePath: "/vela/secret/vault/token",
		Name:     "secret.vault.token",
		Usage:    "token used to access vault system",
	},
	&cli.StringFlag{
		EnvVars:  []string{"VELA_SECRET_VAULT_VERSION", "SECRET_VAULT_VERSION"},
		FilePath: "/vela/secret/vault/version",
		Name:     "secret.vault.version",
		Usage:    "version for the kv backend for the vault system",
		Value:    "2",
	},
}
