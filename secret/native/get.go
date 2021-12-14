// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"github.com/go-vela/types/library"
)

// Get captures a secret.
func (c *client) Get(sType, org, name, path string) (*library.Secret, error) {
	c.Logger.Tracef("getting native %s secret %s for %s/%s", sType, path, org, name)

	// capture the secret from the native service
	s, err := c.Database.GetSecret(sType, org, name, path)
	if err != nil {
		return nil, err
	}

	return s, nil
}
