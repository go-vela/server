// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

//nolint:dupl // ignore similar code with create.go
package secret

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// UpdateSecret updates an existing secret in the database.
func (e *engine) UpdateSecret(s *library.Secret) error {
	// handle the secret based off the type
	switch s.GetType() {
	case constants.SecretShared:
		e.logger.WithFields(logrus.Fields{
			"org":    s.GetOrg(),
			"team":   s.GetTeam(),
			"secret": s.GetName(),
			"type":   s.GetType(),
		}).Tracef("updating secret %s/%s/%s/%s in the database", s.GetType(), s.GetOrg(), s.GetTeam(), s.GetName())
	default:
		e.logger.WithFields(logrus.Fields{
			"org":    s.GetOrg(),
			"repo":   s.GetRepo(),
			"secret": s.GetName(),
			"type":   s.GetType(),
		}).Tracef("updating secret %s/%s/%s/%s in the database", s.GetType(), s.GetOrg(), s.GetRepo(), s.GetName())
	}

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#SecretFromLibrary
	secret := database.SecretFromLibrary(s)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Secret.Validate
	err := secret.Validate()
	if err != nil {
		return err
	}

	// encrypt the fields for the secret
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Secret.Encrypt
	err = secret.Encrypt(e.config.EncryptionKey)
	if err != nil {
		switch s.GetType() {
		case constants.SecretShared:
			return fmt.Errorf("unable to encrypt secret %s/%s/%s/%s: %w", s.GetType(), s.GetOrg(), s.GetTeam(), s.GetName(), err)
		default:
			return fmt.Errorf("unable to encrypt secret %s/%s/%s/%s: %w", s.GetType(), s.GetOrg(), s.GetRepo(), s.GetName(), err)
		}
	}

	// send query to the database
	return e.client.
		Table(constants.TableSecret).
		Save(secret.Nullify()).
		Error
}
