// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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

	// TODO: A bug has been found with this functionality.
	// That affects secret values under 12 chars

	// // encrypt secret value
	// value, err := encrypt([]byte(s.GetValue()), c.passphrase)
	// if err != nil {
	// 	return err
	// }

	// // update value of secret to be encrypted
	// s.SetValue(value)

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
