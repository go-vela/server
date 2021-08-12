// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// List captures a list of secrets.
//
// nolint: lll // ignore long line length
func (c *client) List(sType, org, name string, page, perPage int, teams []string) ([]*library.Secret, error) {
	logrus.Tracef("Listing native %s secrets for %s/%s", sType, org, name)

	// capture the list of secrets from the native service
	s, err := c.Database.GetTypeSecretList(sType, org, name, page, perPage, teams)
	if err != nil {
		return nil, err
	}

	return s, nil
}
