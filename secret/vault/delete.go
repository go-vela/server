// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"fmt"

	"github.com/go-vela/types/constants"
)

// Delete deletes a secret.
func (c *client) Delete(sType, org, name, path string) error {
	c.Logger.Tracef("deleting vault %s secret %s for %s/%s", sType, path, org, name)

	// delete the secret from the Vault service
	switch sType {
	case constants.SecretOrg:
		return c.deleteOrg(org, path)
	case constants.SecretRepo:
		return c.deleteRepo(org, name, path)
	case constants.SecretShared:
		return c.deleteShared(org, name, path)
	default:
		return fmt.Errorf("invalid secret type: %v", sType)
	}
}

// deleteOrg is a helper function to delete
// the org secret for the provided path.
func (c *client) deleteOrg(org, path string) error {
	return c.delete(fmt.Sprintf("%s/org/%s/%s", c.config.Prefix, org, path))
}

// deleteRepo is a helper function to delete
// the repo secret for the provided path.
func (c *client) deleteRepo(org, repo, path string) error {
	return c.delete(fmt.Sprintf("%s/repo/%s/%s/%s", c.config.Prefix, org, repo, path))
}

// deleteShared is a helper function to delete
// the shared secret for the provided path.
func (c *client) deleteShared(org, team, path string) error {
	return c.delete(fmt.Sprintf("%s/shared/%s/%s/%s", c.config.Prefix, org, team, path))
}

// delete is a helper function to delete
// the secret for the provided path.
func (c *client) delete(path string) error {
	// send API call to delete the secret
	_, err := c.Vault.Logical().Delete(path)
	if err != nil {
		return err
	}

	return nil
}
