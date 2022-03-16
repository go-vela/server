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

// Update updates an existing secret.
func (c *client) Update(sType, org, name string, s *library.Secret) error {
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

	c.Logger.WithFields(fields).Tracef("updating native %s secret %s for %s/%s", sType, s.GetName(), org, name)

	// capture the secret from the native service
	sec, err := c.Database.GetSecret(sType, org, name, s.GetName())
	if err != nil {
		return err
	}

	// update the events if set
	if len(s.GetEvents()) > 0 {
		sec.SetEvents(s.GetEvents())
	}

	// update the images if set
	if s.Images != nil {
		sec.SetImages(s.GetImages())
	}

	// update the value if set
	if len(s.GetValue()) > 0 {
		sec.SetValue(s.GetValue())
	}

	// update allow_command if set
	if s.AllowCommand != nil {
		sec.SetAllowCommand(s.GetAllowCommand())
	}

	// update updated_at if set
	sec.SetUpdatedAt(s.GetUpdatedAt())

	// update updated_by if set
	sec.SetUpdatedBy(s.GetUpdatedBy())

	return c.Database.UpdateSecret(sec)
}
