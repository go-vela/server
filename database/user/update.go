// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code in create.go
package user

import (
	"context"
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// UpdateUser updates an existing user in the database.
func (e *engine) UpdateUser(ctx context.Context, u *library.User) (*library.User, error) {
	e.logger.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Tracef("updating user %s in the database", u.GetName())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#UserFromLibrary
	user := database.UserFromLibrary(u)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#User.Validate
	err := user.Validate()
	if err != nil {
		return nil, err
	}

	// encrypt the fields for the user
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#User.Encrypt
	err = user.Encrypt(e.config.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt user %s: %w", u.GetName(), err)
	}

	// send query to the database
	result := e.client.Table(constants.TableUser).Save(user)

	// decrypt fields to return user
	err = user.Decrypt(e.config.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("unable to decrypt user %s: %w", u.GetName(), err)
	}

	return user.ToLibrary(), result.Error
}
