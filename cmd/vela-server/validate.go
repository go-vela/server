// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"github.com/go-vela/types/constants"
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

	if len(c.String("clone-image")) == 0 {
		return fmt.Errorf("clone-image (VELA_CLONE_IMAGE) flag is not properly configured")
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

		if len(c.String("webui-oauth-callback")) == 0 {
			return fmt.Errorf("webui-oauth (VELA_WEBUI_OAUTH_CALLBACK_PATH or VELA_WEBUI_OAUTH_CALLBACK) not set")
		}
	}

	if c.Duration("refresh-token-duration").Seconds() <= c.Duration("access-token-duration").Seconds() {
		return fmt.Errorf("refresh-token-duration (VELA_REFRESH_TOKEN_DURATION) must be larger than the access-token-duration (VELA_ACCESS_TOKEN_DURATION)")
	}

	if c.Int64("default-build-limit") == 0 {
		return fmt.Errorf("default-build-limit (VELA_DEFAULT_BUILD_LIMIT) flag must be greater than 0")
	}

	if c.Int64("max-build-limit") == 0 {
		return fmt.Errorf("max-build-limit (VELA_MAX_BUILD_LIMIT) flag must be greater than 0")
	}

	for _, event := range c.StringSlice("default-repo-events") {
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
