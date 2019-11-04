// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// Update updates a secret.
func (c *client) Update(sType, org, name string, s *library.Secret) error {
	logrus.Tracef("Updating vault %s secret %s for %s/%s", sType, s.GetName(), org, name)

	// capture the secret from the Vault service
	sec, err := c.Get(sType, org, name, s.GetName())
	if err != nil {
		return err
	}

	// convert the Vault secret our secret
	vault := vaultFromSecret(sec)
	if len(s.GetEvents()) > 0 {
		vault.Data["events"] = s.GetEvents()
	}

	if len(s.GetImages()) > 0 {
		vault.Data["images"] = s.GetImages()
	}

	if len(s.GetValue()) > 0 {
		vault.Data["value"] = s.GetValue()
	}

	// validate the secret
	err = database.SecretFromLibrary(secretFromVault(vault)).Validate()
	if err != nil {
		return err
	}

	// update the secret for the Vault service
	switch sType {
	case constants.SecretOrg:
		return c.updateOrg(org, s.GetName(), vault.Data)
	case constants.SecretRepo:
		return c.updateRepo(org, name, s.GetName(), vault.Data)
	case constants.SecretShared:
		fallthrough
	default:
		return c.updateShared(org, name, s.GetName(), vault.Data)
	}
}

// updateOrg is a helper function to update
// the org secret for the provided path.
func (c *client) updateOrg(org, path string, data map[string]interface{}) error {
	return c.update(fmt.Sprintf("secret/%s/%s/%s", constants.SecretOrg, org, path), data)
}

// updateRepo is a helper function to update
// the repo secret for the provided path.
func (c *client) updateRepo(org, repo, path string, data map[string]interface{}) error {
	return c.update(fmt.Sprintf("secret/%s/%s/%s/%s", constants.SecretRepo, org, repo, path), data)
}

// updateShared is a helper function to update
// the shared secret for the provided path.
func (c *client) updateShared(org, team, path string, data map[string]interface{}) error {
	return c.update(fmt.Sprintf("secret/%s/%s/%s/%s", constants.SecretShared, org, team, path), data)
}

// update is a helper function to update
// the secret for the provided path.
func (c *client) update(path string, data map[string]interface{}) error {
	_, err := c.Vault.Logical().Write(path, data)
	if err != nil {
		return err
	}

	return nil
}
