// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// Update updates an existing secret.
func (c *client) Update(sType, org, name string, s *library.Secret) (*library.Secret, error) {
	// capture the secret from the native service
	secret, err := c.Get(sType, org, name, s.GetName())
	if err != nil {
		return nil, err
	}

	// update the events if set
	if len(s.GetEvents()) > 0 {
		secret.SetEvents(s.GetEvents())
	}

	// update the images if set
	if s.Images != nil {
		secret.SetImages(s.GetImages())
	}

	// update the value if set
	if len(s.GetValue()) > 0 {
		secret.SetValue(s.GetValue())
	}

	// update allow_command if set
	if s.AllowCommand != nil {
		secret.SetAllowCommand(s.GetAllowCommand())
	}

	// update updated_at if set
	secret.SetUpdatedAt(s.GetUpdatedAt())

	// update updated_by if set
	secret.SetUpdatedBy(s.GetUpdatedBy())

	// handle the secret based off the type
	switch sType {
	case constants.SecretOrg:
		c.Logger.WithFields(logrus.Fields{
			"org":    org,
			"secret": s.GetName(),
			"type":   sType,
		}).Tracef("updating native %s secret %s for %s", sType, s.GetName(), org)

		// update the org secret in the native service
		return c.Database.UpdateSecret(secret)
	case constants.SecretRepo:
		c.Logger.WithFields(logrus.Fields{
			"org":    org,
			"repo":   name,
			"secret": s.GetName(),
			"type":   sType,
		}).Tracef("updating native %s secret %s for %s/%s", sType, s.GetName(), org, name)

		// update the repo secret in the native service
		return c.Database.UpdateSecret(secret)
	case constants.SecretShared:
		c.Logger.WithFields(logrus.Fields{
			"org":    org,
			"team":   name,
			"secret": s.GetName(),
			"type":   sType,
		}).Tracef("updating native %s secret %s for %s/%s", sType, s.GetName(), org, name)

		// update the shared secret in the native service
		return c.Database.UpdateSecret(secret)
	default:
		return nil, fmt.Errorf("invalid secret type: %s", sType)
	}
}
