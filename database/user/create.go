// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CreateUser creates a new user in the database.
func (e *Engine) CreateUser(ctx context.Context, u *api.User) (*api.User, error) {
	e.logger.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Tracef("creating user %s", u.GetName())

	// cast the API type to database type
	user := types.UserFromAPI(u)

	// validate the necessary fields are populated
	err := user.Validate()
	if err != nil {
		return nil, err
	}

	// encrypt the fields for the user
	err = user.Encrypt(e.config.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt user %s: %w", u.GetName(), err)
	}

	// send query to the database
	result := e.client.
		WithContext(ctx).
		Table(constants.TableUser).
		Create(user)

	// decrypt fields to return user
	err = user.Decrypt(e.config.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("unable to decrypt user %s: %w", u.GetName(), err)
	}

	return user.ToAPI(), result.Error
}
