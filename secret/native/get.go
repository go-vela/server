// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// Get captures a secret.
func (c *client) Get(sType, org, name, path string) (*library.Secret, error) {
	logrus.Tracef("Getting native %s secret %s for %s/%s", sType, path, org, name)

	// capture the secret from the native service
	s, err := c.Native.GetSecret(sType, org, name, path)
	if err != nil {
		return nil, err
	}

	return s, nil
}
