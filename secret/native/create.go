// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

// Create creates a new secret.
func (c *client) Create(sType, org, name string, s *library.Secret) error {
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

	// nolint: lll // ignore long line length due to parameters
	c.Logger.WithFields(fields).Tracef("creating native %s secret %s for %s/%s", sType, s.GetName(), org, name)

	// create the secret for the native service
	switch sType {
	case constants.SecretOrg:
		fallthrough
	case constants.SecretRepo:
		fallthrough
	case constants.SecretShared:
		return c.Database.CreateSecret(s)
	default:
		return fmt.Errorf("invalid secret type: %v", sType)
	}
}
