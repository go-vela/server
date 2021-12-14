// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"fmt"
	"strings"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	"github.com/hashicorp/vault/api"
)

// List captures a list of secrets.
// TODO: Implement fake pagination?
// We drop page and perPage as we are always returning all results.
// Vault API doesn't seem to support pagination. Might result in undesired
// behavior for fetching Vault secrets in paginated manner.
func (c *client) List(sType, org, name string, _, _ int, _ []string) ([]*library.Secret, error) {
	c.Logger.Tracef("listing vault %s secrets for %s/%s", sType, org, name)

	var err error

	s := []*library.Secret{}
	// nolint: staticcheck // ignore false positive
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
		sec, err := c.Get(sType, org, name, key)
		if err != nil {
			return nil, err
		}

		s = append(s, sec)
	}

	return s, nil
}

// listOrg is a helper function to capture the
// list of org secrets for the provided path.
func (c *client) listOrg(org string) (*api.Secret, error) {
	return c.list(fmt.Sprintf("%s/%s/%s", c.config.Prefix, constants.SecretOrg, org))
}

// listRepo is a helper function to capture the
// list of repo secrets for the provided path.
func (c *client) listRepo(org, repo string) (*api.Secret, error) {
	return c.list(fmt.Sprintf("%s/%s/%s/%s", c.config.Prefix, constants.SecretRepo, org, repo))
}

// listShared is a helper function to capture the
// list of shared secrets for the provided path.
func (c *client) listShared(org, team string) (*api.Secret, error) {
	return c.list(fmt.Sprintf("%s/%s/%s/%s", c.config.Prefix, constants.SecretShared, org, team))
}

// list is a helper function to capture the
// list of secrets for the provided path.
func (c *client) list(path string) (*api.Secret, error) {
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
