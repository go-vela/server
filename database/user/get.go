// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package user

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetUser gets a user by ID from the database.
func (e *engine) GetUser(id int64) (*library.User, error) {
	e.logger.Tracef("getting user %d from the database", id)

	// variable to store query results
	u := new(database.User)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableUser).
		Where("id = ?", id).
		Take(u).
		Error
	if err != nil {
		return nil, err
	}

	// decrypt the fields for the user
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#User.Decrypt
	err = u.Decrypt(e.config.EncryptionKey)
	if err != nil {
		// TODO: remove backwards compatibility before 1.x.x release
		//
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allows us to fetch unencrypted users
		e.logger.Errorf("unable to decrypt user %d: %v", u.ID.Int64, err)
	}

	// return the decrypted user
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#User.ToLibrary
	return u.ToLibrary(), nil
}
