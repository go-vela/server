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

// List captures a list of secrets.
func (c *client) List(sType, org, name string, page, perPage int, teams []string) ([]*library.Secret, error) {
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

	c.Logger.WithFields(fields).Tracef("listing native %s secrets for %s/%s", sType, org, name)

	// capture the list of secrets from the native service
	s, err := c.Database.GetTypeSecretList(sType, org, name, page, perPage, teams)
	if err != nil {
		return nil, err
	}

	return s, nil
}
