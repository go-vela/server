// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// Create creates a new secret.
func (c *client) Create(sType, org, name string, s *library.Secret) error {
	logrus.Tracef("Creating native %s secret %s for %s/%s", sType, s.GetName(), org, name)

	// encrypt secret value
	value, err := encrypt([]byte(s.GetValue()), c.passphrase)
	if err != nil {
		return err
	}

	// update value of secret to be encrypted
	s.Value = &value

	// create the secret for the native service
	switch sType {
	case constants.SecretOrg:
		fallthrough
	case constants.SecretRepo:
		fallthrough
	case constants.SecretShared:
		return c.Native.CreateSecret(s)
	default:
		return fmt.Errorf("invalid secret type: %v", sType)
	}
}
