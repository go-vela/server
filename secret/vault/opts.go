// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// ClientOpt represents a configuration option to initialize the secret client.
type ClientOpt func(*client) error

// WithAddress sets the Vault address in the secret client.
func WithAddress(address string) ClientOpt {
	logrus.Trace("configuring address in vault secret client")

	return func(c *client) error {
		// check if the Vault address provided is empty
		if len(address) == 0 {
			return fmt.Errorf("no Vault address provided")
		}

		// set the address in the vault client
		c.config.Address = address

		return nil
	}
}

// WithAuthMethod sets the Vault authentication method in the secret client.
func WithAuthMethod(authMethod string) ClientOpt {
	logrus.Trace("configuring authentication method in vault secret client")

	return func(c *client) error {
		// set the authentication method in the vault client
		c.config.AuthMethod = authMethod

		return nil
	}
}

// WithAWSRole sets the Vault AWS role in the secret client.
func WithAWSRole(awsRole string) ClientOpt {
	logrus.Trace("configuring AWS role in vault secret client")

	return func(c *client) error {
		// set the AWS role in the vault client
		c.config.AWSRole = awsRole

		return nil
	}
}

// WithPrefix sets the Vault prefix in the secret client.
func WithPrefix(prefix string) ClientOpt {
	logrus.Trace("configuring prefix in vault secret client")

	return func(c *client) error {
		// set the prefix in the vault client
		c.config.Prefix = prefix

		return nil
	}
}

// WithToken sets the Vault token in the secret client.
func WithToken(token string) ClientOpt {
	logrus.Trace("configuring token in vault secret client")

	return func(c *client) error {
		// set the token in the vault client
		c.config.Token = token

		return nil
	}
}

// WithTokenDuration sets the Vault token duration in the secret client.
func WithTokenDuration(tokenDuration time.Duration) ClientOpt {
	logrus.Trace("configuring token duration in vault secret client")

	return func(c *client) error {
		// set the token duration in the vault client
		c.config.TokenDuration = tokenDuration

		return nil
	}
}

// WithVersion sets the Vault version in the secret client.
func WithVersion(version string) ClientOpt {
	logrus.Trace("configuring version in vault secret client")

	return func(c *client) error {
		// check if the Vault version provided is empty
		if len(version) == 0 {
			return fmt.Errorf("no Vault version provided")
		}

		// process the vault version being provided
		switch version {
		case "1":
			c.config.SystemPrefix = PrefixVaultV1
		case "2":
			c.config.SystemPrefix = PrefixVaultV2
		default:
			return fmt.Errorf("invalid vault version provided: %s", version)
		}

		// set the version in the vault client
		c.config.Version = version

		return nil
	}
}
