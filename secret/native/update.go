// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// Update updates an existing secret.
func (c *client) Update(sType, org, name string, s *library.Secret) error {
	logrus.Tracef("Updating native %s secret %s for %s/%s", sType, s.GetName(), org, name)

	// capture the secret from the native service
	sec, err := c.Native.GetSecret(sType, org, name, s.GetName())
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

		// encrypt secret value
		value, err := encrypt([]byte(sec.GetValue()), c.passphrase)
		if err != nil {
			return err
		}

		// update value of secret to be encrypted
		sec.Value = &value
	}

	// update allow_command if set
	if s.AllowCommand != nil {
		sec.SetAllowCommand(s.GetAllowCommand())
	}

	return c.Native.UpdateSecret(sec)
}
