// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"strings"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// Get captures a secret.
func (c *client) Get(sType, org, name, path string) (*library.Secret, error) {
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

	c.Logger.WithFields(fields).Tracef("getting native %s secret %s for %s/%s", sType, path, org, name)

	// capture the secret from the native service
	s, err := c.Database.GetSecret(sType, org, name, path)
	if err != nil {
		return nil, err
	}

	return s, nil
}
