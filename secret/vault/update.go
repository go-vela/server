// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// Update updates a secret.
func (c *client) Update(sType, org, name string, s *library.Secret) (*library.Secret, error) {
	// create log fields from secret metadata
	fields := logrus.Fields{
		"org":    org,
		"repo":   name,
		"secret": s.GetName(),
		"type":   sType,
	}

	// check if secret is a shared secret
	if strings.EqualFold(sType, constants.SecretShared) {
		// update log fields from secret metadata
		fields = logrus.Fields{
			"org":    org,
			"team":   name,
			"secret": s.GetName(),
			"type":   sType,
		}
	}

	c.Logger.WithFields(fields).Tracef("updating vault %s secret %s for %s/%s", sType, s.GetName(), org, name)

	// capture the secret from the Vault service
	sec, err := c.Get(sType, org, name, s.GetName())
	if err != nil {
		return nil, err
	}

	// convert the Vault secret our secret
	vault := vaultFromSecret(sec)
	if len(s.GetEvents()) > 0 {
		vault.Data["events"] = s.GetEvents()
	}

	if s.Images != nil {
		vault.Data["images"] = s.GetImages()
	}

	if len(s.GetValue()) > 0 {
		vault.Data["value"] = s.GetValue()
	}

	if s.AllowCommand != nil {
		vault.Data["allow_command"] = s.GetAllowCommand()
	}

	// validate the secret
	err = database.SecretFromLibrary(secretFromVault(vault)).Validate()
	if err != nil {
		return nil, err
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
func (c *client) updateOrg(org, path string, data map[string]interface{}) (*library.Secret, error) {
	return c.update(fmt.Sprintf("%s/%s/%s/%s", c.config.Prefix, constants.SecretOrg, org, path), data)
}

// updateRepo is a helper function to update
// the repo secret for the provided path.
func (c *client) updateRepo(org, repo, path string, data map[string]interface{}) (*library.Secret, error) {
	return c.update(fmt.Sprintf("%s/%s/%s/%s/%s", c.config.Prefix, constants.SecretRepo, org, repo, path), data)
}

// updateShared is a helper function to update
// the shared secret for the provided path.
func (c *client) updateShared(org, team, path string, data map[string]interface{}) (*library.Secret, error) {
	return c.update(fmt.Sprintf("%s/%s/%s/%s/%s", c.config.Prefix, constants.SecretShared, org, team, path), data)
}

// update is a helper function to update
// the secret for the provided path.
func (c *client) update(path string, data map[string]interface{}) (*library.Secret, error) {
	if strings.HasPrefix("secret/data", c.config.Prefix) {
		data = map[string]interface{}{
			"data": data,
		}
	}

	s, err := c.Vault.Logical().Write(path, data)
	if err != nil {
		return nil, err
	}

	return secretFromVault(s), nil
}
