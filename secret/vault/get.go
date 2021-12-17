// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/hashicorp/vault/api"
)

// Get captures a secret.
func (c *client) Get(sType, org, name, path string) (s *library.Secret, err error) {
	// create log fields from secret metadata
	fields := logrus.Fields{
		"org":    org,
		"repo":   name,
		"secret": path,
		"type":   sType,
	}

	// check if secret is a shared secret
	if strings.EqualFold(sType, constants.SecretShared) {
		// update log fields from secret metadata
		fields = logrus.Fields{
			"org":    org,
			"team":   name,
			"secret": path,
			"type":   sType,
		}
	}

	c.Logger.WithFields(fields).Tracef("getting vault %s secret %s for %s/%s", sType, path, org, name)

	// nolint: ineffassign,staticcheck // ignore false positive
	vault := new(api.Secret)

	// capture the secret from the Vault service
	switch sType {
	case constants.SecretOrg:
		vault, err = c.getOrg(org, path)
	case constants.SecretRepo:
		vault, err = c.getRepo(org, name, path)
	case constants.SecretShared:
		vault, err = c.getShared(org, name, path)
	default:
		return nil, fmt.Errorf("invalid secret type: %v", sType)
	}

	if err != nil {
		return nil, err
	}

	return secretFromVault(vault), nil
}

// getOrg is a helper function to capture
// the org secret for the provided path.
func (c *client) getOrg(org, path string) (*api.Secret, error) {
	return c.get(fmt.Sprintf("%s/%s/%s/%s", c.config.Prefix, constants.SecretOrg, org, path))
}

// getRepo is a helper function to capture
// the repo secret for the provided path.
func (c *client) getRepo(org, repo, path string) (*api.Secret, error) {
	return c.get(fmt.Sprintf("%s/%s/%s/%s/%s", c.config.Prefix, constants.SecretRepo, org, repo, path))
}

// getShared is a helper function to capture
// the shared secret for the provided path.
func (c *client) getShared(org, team, path string) (*api.Secret, error) {
	// nolint: lll // ignore long line length due to parameters
	return c.get(fmt.Sprintf("%s/%s/%s/%s/%s", c.config.Prefix, constants.SecretShared, org, team, path))
}

// get is a helper function to capture
// the secret for the provided path.
func (c *client) get(path string) (*api.Secret, error) {
	// send API call to capture the secret
	vault, err := c.Vault.Logical().Read(path)
	if err != nil {
		return nil, err
	}

	// return nil if secret does not exist
	if vault == nil {
		return nil, fmt.Errorf("secret %s does not exist", path)
	}

	return vault, nil
}
