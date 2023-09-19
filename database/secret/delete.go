// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// DeleteSecret deletes an existing secret from the database.
func (e *engine) DeleteSecret(ctx context.Context, s *library.Secret) error {
	// handle the secret based off the type
	//
	//nolint:dupl // ignore similar code with update.go
	switch s.GetType() {
	case constants.SecretShared:
		e.logger.WithFields(logrus.Fields{
			"org":    s.GetOrg(),
			"team":   s.GetTeam(),
			"secret": s.GetName(),
			"type":   s.GetType(),
		}).Tracef("deleting secret %s/%s/%s/%s from the database", s.GetType(), s.GetOrg(), s.GetTeam(), s.GetName())
	default:
		e.logger.WithFields(logrus.Fields{
			"org":    s.GetOrg(),
			"repo":   s.GetRepo(),
			"secret": s.GetName(),
			"type":   s.GetType(),
		}).Tracef("deleting secret %s/%s/%s/%s from the database", s.GetType(), s.GetOrg(), s.GetRepo(), s.GetName())
	}

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#SecretFromLibrary
	secret := database.SecretFromLibrary(s)

	// send query to the database
	return e.client.
		Table(constants.TableSecret).
		Delete(secret).
		Error
}
