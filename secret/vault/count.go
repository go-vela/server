// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"fmt"

	"github.com/go-vela/types/constants"

	"github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
)

// Count counts a list of secrets.
func (c *client) Count(sType, org, name string) (i int64, err error) {
	logrus.Tracef("Counting vault %s secrets for %s/%s", sType, org, name)

	// nolint: staticcheck // ignore false positive
	vault := new(api.Secret)
	count := 0

	// capture the list of secrets from the Vault service
	switch sType {
	case constants.SecretOrg:
		vault, err = c.listOrg(org)
	case constants.SecretRepo:
		vault, err = c.listRepo(org, name)
	case constants.SecretShared:
		vault, err = c.listShared(org, name)
	default:
		return 0, fmt.Errorf("invalid secret type: %v", sType)
	}

	if err != nil {
		return 0, err
	}

	// cast the list of secrets to the expected type
	keys, ok := vault.Data["keys"].([]interface{})
	if !ok {
		return 0, fmt.Errorf("not a valid list of secrets from Vault")
	}

	// iterate through each element in the list of secrets
	for _, element := range keys {
		// cast the secret to the expected type
		_, ok := element.(string)
		if !ok {
			return 0, fmt.Errorf("not a valid list of secrets from Vault")
		}

		count++
	}

	return int64(count), nil
}
