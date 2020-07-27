// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func validate(c *cli.Context) error {
	logrus.Debug("Validating CLI configuration")

	// validate core configuration
	err := validateCore(c)
	if err != nil {
		return err
	}

	// validate compiler configuration
	err = validateCompiler(c)
	if err != nil {
		return err
	}

	// validate database configuration
	err = validateDatabase(c)
	if err != nil {
		return err
	}

	// validate queue configuration
	err = validateQueue(c)
	if err != nil {
		return err
	}

	// validate secret configuration
	err = validateSecret(c)
	if err != nil {
		return err
	}

	// validate source configuration
	err = validateSource(c)
	if err != nil {
		return err
	}

	return nil
}

// helper function to validate the core CLI configuration.
func validateCore(c *cli.Context) error {
	logrus.Trace("Validating core CLI configuration")

	if len(c.String("server-addr")) == 0 {
		return fmt.Errorf("server-addr (VELA_ADDR or VELA_HOST) flag is not properly configured")
	}

	if !strings.Contains(c.String("server-addr"), "://") {
		return fmt.Errorf("server-addr (VELA_ADDR or VELA_HOST) flag must be <scheme>://<hostname> format")
	}

	if strings.HasSuffix(c.String("server-addr"), "/") {
		return fmt.Errorf("server-addr (VELA_ADDR or VELA_HOST) flag must not have trailing slash")
	}

	if len(c.String("vela-secret")) == 0 {
		return fmt.Errorf("vela-secret (VELA_SECRET) flag is not properly configured")
	}

	if len(c.String("webui-addr")) == 0 {
		logrus.Warn("optional flag webui-addr (VELA_WEBUI_ADDR or VELA_WEBUI_HOST) not set")
	} else {
		if !strings.Contains(c.String("webui-addr"), "://") {
			return fmt.Errorf("webui-addr (VELA_WEBUI_ADDR or VELA_WEBUI_HOST) flag must be <scheme>://<hostname> format")
		}

		if strings.HasSuffix(c.String("webui-addr"), "/") {
			return fmt.Errorf("webui-addr (VELA_WEBUI_ADDR or VELA_WEBUI_HOST) flag must not have trailing slash")
		}
	}

	if c.Duration("refresh-token-duration").Seconds() <= c.Duration("access-token-duration").Seconds() {
		return fmt.Errorf("refresh-token-duration (VELA_REFRESH_TOKEN_DURATION) must be larger than the access-token-duration (VELA_ACCESS_TOKEN_DURATION)")
	}

	return nil
}

// helper function to validate the compiler CLI configuration.
func validateCompiler(c *cli.Context) error {
	logrus.Trace("Validating compiler CLI configuration")

	if c.Bool("github-driver") {
		if len(c.String("github-url")) == 0 {
			return fmt.Errorf("github-url (VELA_COMPILER_GITHUB_URL or COMPILER_GITHUB_URL) flag not specified")
		}

		if len(c.String("github-token")) == 0 {
			return fmt.Errorf("github-token (VELA_COMPILER_GITHUB_TOKEN or COMPILER_GITHUB_TOKEN) flag not specified")
		}
	}

	return nil
}

// helper function to validate the database CLI configuration.
func validateDatabase(c *cli.Context) error {
	logrus.Trace("Validating database CLI configuration")

	if len(c.String("database.driver")) == 0 {
		return fmt.Errorf("database.driver (VELA_DATABASE_DRIVER or DATABASE_DRIVER) flag not specified")
	}

	if len(c.String("database.config")) == 0 {
		return fmt.Errorf("database.config (VELA_DATABASE_CONFIG or DATABASE_CONFIG) flag not specified")
	}

	return nil
}

// helper function to validate the queue CLI configuration.
func validateQueue(c *cli.Context) error {
	logrus.Trace("Validating queue CLI configuration")

	if len(c.String("queue-driver")) == 0 {
		return fmt.Errorf("queue-driver (VELA_QUEUE_DRIVER or QUEUE_DRIVER) flag not specified")
	}

	if len(c.String("queue-config")) == 0 {
		return fmt.Errorf("queue-config (VELA_QUEUE_CONFIG or QUEUE_CONFIG) flag not specified")
	}

	return nil
}

// helper function to validate the secret CLI configuration.
func validateSecret(c *cli.Context) error {
	logrus.Trace("Validating secret CLI configuration")

	if c.Bool("vault-driver") {
		if len(c.String("vault-addr")) == 0 {
			return fmt.Errorf("vault-addr (VELA_SECRET_VAULT_ADDR or SECRET_VAULT_ADDR) flag not specified")
		}

		if len(c.String("vault-token")) == 0 {
			return fmt.Errorf("vault-token (VELA_SECRET_VAULT_TOKEN or SECRET_VAULT_TOKEN) flag not specified")
		}
	}

	return nil
}

// helper function to validate the source CLI configuration.
func validateSource(c *cli.Context) error {
	logrus.Trace("Validating source CLI configuration")

	if len(c.String("source-driver")) > 0 {
		if len(c.String("source-url")) == 0 {
			return fmt.Errorf("source-url (VELA_SOURCE_URL or SOURCE_URL) flag not specified")
		}

		if len(c.String("source-client")) == 0 {
			return fmt.Errorf("source-client (VELA_SOURCE_CLIENT or SOURCE_CLIENT) flag not specified")
		}

		if len(c.String("source-secret")) == 0 {
			return fmt.Errorf("source-secret (VELA_SOURCE_SECRET or SOURCE_SECRET) flag not specified")
		}
	}

	return nil
}
