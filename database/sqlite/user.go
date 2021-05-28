// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"errors"
	"fmt"

	"github.com/go-vela/server/database/sqlite/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
)

// GetUser gets a user by unique ID from the database.
//
// nolint: dupl // ignore false positive of duplicate code
func (c *client) GetUser(id int64) (*library.User, error) {
	logrus.Tracef("getting user %d from the database", id)

	// variable to store query results
	u := new(database.User)

	// send query to the database and store result in variable
	result := c.Sqlite.
		Table(constants.TableUser).
		Raw(dml.SelectUser, id).
		Scan(u)

	// check if the query returned a record not found error or no rows were returned
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// decrypt the fields for the user
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#User.Decrypt
	err := u.Decrypt(c.config.EncryptionKey)
	if err != nil {
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allows us to fetch unencrypted users
		logrus.Errorf("unable to decrypt user %d: %v", id, err)

		// return the unencrypted user
		return u.ToLibrary(), result.Error
	}

	// return the decrypted user
	return u.ToLibrary(), result.Error
}

// GetUserName gets a user by name from the database.
//
// nolint: dupl // ignore false positive of duplicate code
func (c *client) GetUserName(name string) (*library.User, error) {
	logrus.Tracef("getting user %s from the database", name)

	// variable to store query results
	u := new(database.User)

	// send query to the database and store result in variable
	result := c.Sqlite.
		Table(constants.TableUser).
		Raw(dml.SelectUserName, name).
		Scan(u)

	// check if the query returned a record not found error or no rows were returned
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// decrypt the fields for the user
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#User.Decrypt
	err := u.Decrypt(c.config.EncryptionKey)
	if err != nil {
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allows us to fetch unencrypted users
		logrus.Errorf("unable to decrypt user %s: %v", name, err)

		// return the unencrypted user
		return u.ToLibrary(), result.Error
	}

	// return the decrypted user
	return u.ToLibrary(), result.Error
}

// CreateUser creates a new user in the database.
func (c *client) CreateUser(u *library.User) error {
	logrus.Tracef("creating user %s from the database", u.GetName())

	// cast to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#UserFromLibrary
	user := database.UserFromLibrary(u)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#User.Validate
	err := user.Validate()
	if err != nil {
		return err
	}

	// encrypt the fields for the user
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#User.Encrypt
	err = user.Encrypt(c.config.EncryptionKey)
	if err != nil {
		return fmt.Errorf("unable to encrypt user %s: %v", u.GetName(), err)
	}

	// send query to the database
	return c.Sqlite.
		Table(constants.TableUser).
		Create(user).Error
}

// UpdateUser updates a user in the database.
func (c *client) UpdateUser(u *library.User) error {
	logrus.Tracef("updating user %s from the database", u.GetName())

	// cast to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#UserFromLibrary
	user := database.UserFromLibrary(u)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#User.Validate
	err := user.Validate()
	if err != nil {
		return err
	}

	// encrypt the fields for the user
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#User.Encrypt
	err = user.Encrypt(c.config.EncryptionKey)
	if err != nil {
		return fmt.Errorf("unable to encrypt user %s: %v", u.GetName(), err)
	}

	// send query to the database
	return c.Sqlite.
		Table(constants.TableUser).
		Save(user).Error
}

// DeleteUser deletes a user by unique ID from the database.
func (c *client) DeleteUser(id int64) error {
	logrus.Tracef("deleting user %d from the database", id)

	// send query to the database
	return c.Sqlite.
		Table(constants.TableUser).
		Exec(dml.DeleteUser, id).Error
}
