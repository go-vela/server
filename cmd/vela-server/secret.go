// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/secret"
	"github.com/go-vela/server/secret/native"
	"github.com/go-vela/server/secret/vault"
	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

// helper function to setup the secrets engines from the CLI arguments.
func setupSecrets(c *cli.Context, d database.Service) (map[string]secret.Service, error) {
	logrus.Debug("Creating secret clients from CLI configuration")

	secrets := make(map[string]secret.Service)

	native, err := setupNative(d)
	if err != nil {
		return nil, err
	}

	secrets[constants.DriverNative] = native

	if c.Bool("vault-driver") {
		vault, err := setupVault(c)
		if err != nil {
			return nil, err
		}

		secrets[constants.DriverVault] = vault
	}

	return secrets, nil
}

// helper function to setup the Native secret engine from the CLI arguments.
func setupNative(d database.Service) (secret.Service, error) {
	logrus.Tracef("Creating %s secret client from CLI configuration", constants.DriverNative)
	return native.New(d)
}

// helper function to setup the Vault secret engine from the CLI arguments.
func setupVault(c *cli.Context) (secret.Service, error) {
	logrus.Tracef("Creating %s secret client from CLI configuration", constants.DriverVault)

	// create new Vault secret service
	//
	// https://pkg.go.dev/github.com/go-vela/server/secret/vault?tab=doc#New
	return vault.New(
		vault.WithAddress(c.String("vault-addr")),
		vault.WithAuthMethod(c.String("vault-auth-method")),
		vault.WithAWSRole(c.String("vault-aws-role")),
		vault.WithPrefix(c.String("vault-prefix")),
		vault.WithToken(c.String("vault-token")),
		vault.WithTokenDuration(c.Duration("vault-renewal")),
		vault.WithVersion(c.String("vault-version")),
	)
}
