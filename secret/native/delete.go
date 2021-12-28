// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"strings"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// Delete deletes a secret.
func (c *client) Delete(sType, org, name, path string) error {
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

	// nolint: lll // ignore long line length due to parameters
	c.Logger.WithFields(fields).Tracef("deleting native %s secret %s for %s/%s", sType, path, org, name)

	// capture the secret from the native service
	s, err := c.Database.GetSecret(sType, org, name, path)
	if err != nil {
		return err
	}

	// delete the secret from the native service
	return c.Database.DeleteSecret(s.GetID())
}
