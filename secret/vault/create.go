// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"fmt"
	"strings"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// Create creates a new secret.
func (c *client) Create(sType, org, name string, s *library.Secret) error {
	logrus.Tracef("Creating vault %s secret %s for %s/%s", sType, s.GetName(), org, name)

	// validate the secret
	err := database.SecretFromLibrary(s).Validate()
	if err != nil {
		return err
	}

	// convert our secret to a Vault secret
	vault := vaultFromSecret(s)

	// create the secret for the Vault service
	switch sType {
	case constants.SecretOrg:
		return c.createOrg(org, s.GetName(), vault.Data)
	case constants.SecretRepo:
		return c.createRepo(org, name, s.GetName(), vault.Data)
	case constants.SecretShared:
		return c.createShared(org, name, s.GetName(), vault.Data)
	default:
		return fmt.Errorf("invalid secret type: %v", sType)
	}
}

// createOrg is a helper function to create
// the org secret for the provided path.
func (c *client) createOrg(org, path string, data map[string]interface{}) error {
	return c.create(fmt.Sprintf("%s/org/%s/%s", c.config.Prefix, org, path), data)
}

// createRepo is a helper function to create
// the repo secret for the provided path.
func (c *client) createRepo(org, repo, path string, data map[string]interface{}) error {
	return c.create(fmt.Sprintf("%s/repo/%s/%s/%s", c.config.Prefix, org, repo, path), data)
}

// createShared is a helper function to create
// the shared secret for the provided path.
func (c *client) createShared(org, team, path string, data map[string]interface{}) error {
	return c.create(fmt.Sprintf("%s/shared/%s/%s/%s", c.config.Prefix, org, team, path), data)
}

// create is a helper function to create
// the secret for the provided path.
func (c *client) create(path string, data map[string]interface{}) error {
	if strings.HasPrefix("secret/data", c.config.Prefix) {
		data = map[string]interface{}{
			"data": data,
		}
	}

	// send API call to create the secret
	_, err := c.Vault.Logical().Write(path, data)
	if err != nil {
		return err
	}

	return nil
}
