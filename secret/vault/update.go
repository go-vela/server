// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	database "github.com/go-vela/server/database/types"
)

// Update updates a secret.
func (c *Client) Update(ctx context.Context, sType, org, name string, s *api.Secret) (*api.Secret, error) {
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
	sec, err := c.Get(ctx, sType, org, name, s.GetName())
	if err != nil {
		return nil, err
	}

	// convert the Vault secret our secret
	vault := vaultFromSecret(sec)

	if s.GetAllowEvents().ToDatabase() != 0 {
		vault.Data["allow_events"] = s.GetAllowEvents().ToDatabase()
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

	if s.AllowSubstitution != nil {
		vault.Data["allow_substitution"] = s.GetAllowSubstitution()
	}

	// validate the secret
	err = database.SecretFromAPI(secretFromVault(vault)).Validate()
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
func (c *Client) updateOrg(org, path string, data map[string]any) (*api.Secret, error) {
	return c.update(fmt.Sprintf("%s/%s/%s/%s", c.config.Prefix, constants.SecretOrg, org, path), data)
}

// updateRepo is a helper function to update
// the repo secret for the provided path.
func (c *Client) updateRepo(org, repo, path string, data map[string]any) (*api.Secret, error) {
	return c.update(fmt.Sprintf("%s/%s/%s/%s/%s", c.config.Prefix, constants.SecretRepo, org, repo, path), data)
}

// updateShared is a helper function to update
// the shared secret for the provided path.
func (c *Client) updateShared(org, team, path string, data map[string]any) (*api.Secret, error) {
	return c.update(fmt.Sprintf("%s/%s/%s/%s/%s", c.config.Prefix, constants.SecretShared, org, team, path), data)
}

// update is a helper function to update
// the secret for the provided path.
func (c *Client) update(path string, data map[string]any) (*api.Secret, error) {
	if strings.HasPrefix("secret/data", c.config.Prefix) {
		data = map[string]any{
			"data": data,
		}
	}

	s, err := c.Vault.Logical().Write(path, data)
	if err != nil {
		return nil, err
	}

	return secretFromVault(s), nil
}
