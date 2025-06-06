// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"

	velaAPI "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// List captures a list of secrets.
// TODO: Implement fake pagination?
// We drop page and perPage as we are always returning all results.
// Vault API doesn't seem to support pagination. Might result in undesired
// behavior for fetching Vault secrets in paginated manner.
func (c *Client) List(ctx context.Context, sType, org, name string, _, _ int, _ []string) ([]*velaAPI.Secret, error) {
	// create log fields from secret metadata
	fields := logrus.Fields{
		"org":  org,
		"repo": name,
		"type": sType,
	}

	// check if secret is a shared secret
	if strings.EqualFold(sType, constants.SecretShared) {
		// update log fields from secret metadata
		fields = logrus.Fields{
			"org":  org,
			"team": name,
			"type": sType,
		}
	}

	c.Logger.WithFields(fields).Tracef("listing vault %s secrets for %s/%s", sType, org, name)

	var err error

	s := []*velaAPI.Secret{}

	vault := new(api.Secret)

	// capture the list of secrets from the Vault service
	switch sType {
	case constants.SecretOrg:
		vault, err = c.listOrg(org)
	case constants.SecretRepo:
		vault, err = c.listRepo(org, name)
	case constants.SecretShared:
		vault, err = c.listShared(org, name)
	default:
		return nil, fmt.Errorf("invalid secret type: %v", sType)
	}

	if err != nil {
		return nil, err
	}

	// cast the list of secrets to the expected type
	keys, ok := vault.Data["keys"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid list of secrets from Vault")
	}

	// iterate through each element in the list of secrets
	for _, element := range keys {
		// cast the secret to the expected type
		key, ok := element.(string)
		if !ok {
			return nil, fmt.Errorf("not a valid list of secrets from Vault")
		}

		// capture the secret from the Vault service
		sec, err := c.Get(ctx, sType, org, name, key)
		if err != nil {
			return nil, err
		}

		s = append(s, sec)
	}

	return s, nil
}

// listOrg is a helper function to capture the
// list of org secrets for the provided path.
func (c *Client) listOrg(org string) (*api.Secret, error) {
	return c.list(fmt.Sprintf("%s/%s/%s", c.config.Prefix, constants.SecretOrg, org))
}

// listRepo is a helper function to capture the
// list of repo secrets for the provided path.
func (c *Client) listRepo(org, repo string) (*api.Secret, error) {
	return c.list(fmt.Sprintf("%s/%s/%s/%s", c.config.Prefix, constants.SecretRepo, org, repo))
}

// listShared is a helper function to capture the
// list of shared secrets for the provided path.
func (c *Client) listShared(org, team string) (*api.Secret, error) {
	return c.list(fmt.Sprintf("%s/%s/%s/%s", c.config.Prefix, constants.SecretShared, org, team))
}

// list is a helper function to capture the
// list of secrets for the provided path.
func (c *Client) list(path string) (*api.Secret, error) {
	// handle k/v v2
	if strings.HasPrefix(path, "secret/data/") {
		// remove secret/data/ prefix
		path = strings.TrimPrefix(path, "secret/data/")
		// add secret/metadata/ prefix
		path = fmt.Sprintf("secret/metadata/%s", path)
	}

	// send API call to capture the list of secrets
	vault, err := c.Vault.Logical().List(path)
	if err != nil {
		return nil, err
	}

	return vault, nil
}
