// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
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
		sec.Events = s.Events
	}

	// update the images if set
	if len(s.GetImages()) > 0 {
		sec.Images = s.Images
	}

	// update the value if set
	if len(s.GetValue()) > 0 {
		sec.Value = s.Value
	}

	return c.Native.UpdateSecret(sec)
}
