// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"fmt"
	"time"
)

// ClientOpt represents a configuration option to initialize the secret client for Vault.
type ClientOpt func(*client) error

// WithAddress sets the address in the secret client for Vault.
func WithAddress(address string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring address in vault secret client")

		// check if the Vault address provided is empty
		if len(address) == 0 {
			return fmt.Errorf("no Vault address provided")
		}

		// set the address in the vault client
		c.config.Address = address

		return nil
	}
}

// WithAuthMethod sets the authentication method in the secret client for Vault.
func WithAuthMethod(authMethod string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring authentication method in vault secret client")

		// set the authentication method in the vault client
		c.config.AuthMethod = authMethod

		return nil
	}
}

// WithAWSRole sets the AWS role in the secret client for Vault.
func WithAWSRole(awsRole string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring AWS role in vault secret client")

		// set the AWS role in the vault client
		c.config.AWSRole = awsRole

		return nil
	}
}

// WithPrefix sets the prefix in the secret client for Vault.
func WithPrefix(prefix string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring prefix in vault secret client")

		// set the prefix in the vault client
		c.config.Prefix = prefix

		return nil
	}
}

// WithToken sets the token in the secret client for Vault.
func WithToken(token string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring token in vault secret client")

		// set the token in the vault client
		c.config.Token = token

		return nil
	}
}

// WithTokenDuration sets the token duration in the secret client for Vault.
func WithTokenDuration(tokenDuration time.Duration) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring token duration in vault secret client")

		// set the token duration in the vault client
		c.config.TokenDuration = tokenDuration

		return nil
	}
}

// WithVersion sets the version in the secret client for Vault.
func WithVersion(version string) ClientOpt {
	return func(c *client) error {
		c.Logger.Trace("configuring version in vault secret client")

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
