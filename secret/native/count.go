// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"github.com/sirupsen/logrus"
)

// Count counts a list of secrets.
func (c *client) Count(sType, org, name string, teams []string) (int64, error) {
	logrus.Tracef("Counting native %s secrets for %s/%s", sType, org, name)

	// capture the count of secrets from the native service
	s, err := c.Database.GetTypeSecretCount(sType, org, name, teams)
	if err != nil {
		return 0, err
	}

	return s, nil
}
