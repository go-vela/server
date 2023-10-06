// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/secret"
	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

// helper function to setup the secrets engines from the CLI arguments.
func setupSecrets(c *cli.Context, d database.Interface) (map[string]secret.Service, error) {
	logrus.Debug("Creating secret clients from CLI configuration")

	secrets := make(map[string]secret.Service)

	// native secret configuration
	_native := &secret.Setup{
		Driver:   constants.DriverNative,
		Database: d,
	}

	// setup the native secret service
	//
	// https://pkg.go.dev/github.com/go-vela/server/secret?tab=doc#New
	native, err := secret.New(_native)
	if err != nil {
		return nil, err
	}

	secrets[constants.DriverNative] = native

	// check if the vault driver is enabled
	if c.Bool("secret.vault.driver") {
		// vault secret configuration
		_vault := &secret.Setup{
			Driver:        constants.DriverVault,
			Address:       c.String("secret.vault.addr"),
			AuthMethod:    c.String("secret.vault.auth-method"),
			AwsRole:       c.String("secret.vault.aws-role"),
			Prefix:        c.String("secret.vault.prefix"),
			Token:         c.String("secret.vault.token"),
			TokenDuration: c.Duration("secret.vault.renewal"),
			Version:       c.String("secret.vault.version"),
		}

		// setup the vault secret service
		//
		// https://pkg.go.dev/github.com/go-vela/server/secret?tab=doc#New
		vault, err := secret.New(_vault)
		if err != nil {
			return nil, err
		}

		secrets[constants.DriverVault] = vault
	}

	return secrets, nil
}
