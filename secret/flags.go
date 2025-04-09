// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"time"

	"github.com/urfave/cli/v3"
)

// Flags represents all supported command line
// interface (CLI) flags for the secret.
//
// https://pkg.go.dev/github.com/urfave/cli?tab=doc#Flag
var Flags = []cli.Flag{
	// Secret Flags

	&cli.BoolFlag{
		Name:  "secret.vault.driver",
		Usage: "enables the vault secret driver",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SECRET_VAULT"),
			cli.EnvVar("SECRET_VAULT"),
			cli.File("/vela/secret/vault/driver"),
		),
	},
	&cli.StringFlag{
		Name:  "secret.vault.addr",
		Usage: "fully qualified url (<scheme>://<host>) for the vault system",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SECRET_VAULT_ADDR"),
			cli.EnvVar("SECRET_VAULT_ADDR"),
			cli.File("/vela/secret/vault/addr"),
		),
	},
	&cli.StringFlag{
		Name:  "secret.vault.auth-method",
		Usage: "authentication method used to obtain token from vault system",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SECRET_VAULT_AUTH_METHOD"),
			cli.EnvVar("SECRET_VAULT_AUTH_METHOD"),
			cli.File("/vela/secret/vault/auth_method"),
		),
	},
	&cli.StringFlag{
		Name:  "secret.vault.aws-role",
		Usage: "vault role used to connect to the auth/aws/login endpoint",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SECRET_VAULT_AWS_ROLE"),
			cli.EnvVar("SECRET_VAULT_AWS_ROLE"),
			cli.File("/vela/secret/vault/aws_role"),
		),
	},
	&cli.StringFlag{
		Name:  "secret.vault.prefix",
		Usage: "prefix for k/v secrets in vault system e.g. secret/data/<prefix>/<path>",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SECRET_VAULT_PREFIX"),
			cli.EnvVar("SECRET_VAULT_PREFIX"),
			cli.File("/vela/secret/vault/prefix"),
		),
	},
	&cli.DurationFlag{
		Name:  "secret.vault.renewal",
		Usage: "frequency which the vault token should be renewed",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SECRET_VAULT_RENEWAL"),
			cli.EnvVar("SECRET_VAULT_RENEWAL"),
			cli.File("/vela/secret/vault/renewal"),
		),
		Value: 30 * time.Minute,
	},
	&cli.StringFlag{
		Name:  "secret.vault.token",
		Usage: "token used to access vault system",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SECRET_VAULT_TOKEN"),
			cli.EnvVar("SECRET_VAULT_TOKEN"),
			cli.File("/vela/secret/vault/token"),
		),
	},
	&cli.StringFlag{
		Name:  "secret.vault.version",
		Usage: "version for the kv backend for the vault system",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("VELA_SECRET_VAULT_VERSION"),
			cli.EnvVar("SECRET_VAULT_VERSION"),
			cli.File("/vela/secret/vault/version"),
		),
		Value: "2",
	},
}
