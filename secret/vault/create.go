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

// Create creates a new secret.
func (c *Client) Create(_ context.Context, sType, org, name string, s *api.Secret) (*api.Secret, error) {
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

	c.Logger.WithFields(fields).Tracef("creating vault %s secret %s for %s/%s", sType, s.GetName(), org, name)

	// validate the secret
	err := database.SecretFromAPI(s).Validate()
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("invalid secret type: %v", sType)
	}
}

// createOrg is a helper function to create
// the org secret for the provided path.
func (c *Client) createOrg(org, path string, data map[string]interface{}) (*api.Secret, error) {
	return c.create(fmt.Sprintf("%s/org/%s/%s", c.config.Prefix, org, path), data)
}

// createRepo is a helper function to create
// the repo secret for the provided path.
func (c *Client) createRepo(org, repo, path string, data map[string]interface{}) (*api.Secret, error) {
	return c.create(fmt.Sprintf("%s/repo/%s/%s/%s", c.config.Prefix, org, repo, path), data)
}

// createShared is a helper function to create
// the shared secret for the provided path.
func (c *Client) createShared(org, team, path string, data map[string]interface{}) (*api.Secret, error) {
	return c.create(fmt.Sprintf("%s/shared/%s/%s/%s", c.config.Prefix, org, team, path), data)
}

// create is a helper function to create
// the secret for the provided path.
func (c *Client) create(path string, data map[string]interface{}) (*api.Secret, error) {
	if strings.HasPrefix("secret/data", c.config.Prefix) {
		data = map[string]interface{}{
			"data": data,
		}
	}

	// send API call to create the secret
	s, err := c.Vault.Logical().Write(path, data)
	if err != nil {
		return nil, err
	}

	return secretFromVault(s), nil
}
