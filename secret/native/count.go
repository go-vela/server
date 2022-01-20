// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"strings"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// Count counts a list of secrets.
func (c *client) Count(sType, org, name string, teams []string) (int64, error) {
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

	c.Logger.WithFields(fields).Tracef("counting native %s secrets for %s/%s", sType, org, name)

	// capture the count of secrets from the native service
	s, err := c.Database.GetTypeSecretCount(sType, org, name, teams)
	if err != nil {
		return 0, err
	}

	return s, nil
}
