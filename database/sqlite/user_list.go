// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"github.com/go-vela/server/database/sqlite/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// GetUserList gets a list of all users from the database.
func (c *client) GetUserList() ([]*library.User, error) {
	logrus.Trace("listing users from the database")

	// variable to store query results
	u := new([]database.User)

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableUser).
		Raw(dml.ListUsers).
		Scan(u).Error
	if err != nil {
		return nil, err
	}

	// variable we want to return
	users := []*library.User{}
	// iterate through all query results
	for _, user := range *u {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := user

		// decrypt the fields for the user
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#User.Decrypt
		err = tmp.Decrypt(c.config.EncryptionKey)
		if err != nil {
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch unencrypted users
			logrus.Errorf("unable to decrypt user %d: %v", tmp.ID.Int64, err)
		}

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#User.ToLibrary
		users = append(users, tmp.ToLibrary())
	}

	return users, nil
}

// GetUserLiteList gets a lite list of all users from the database.
func (c *client) GetUserLiteList(page, perPage int) ([]*library.User, error) {
	logrus.Trace("listing lite users from the database")

	// variable to store query results
	u := new([]database.User)
	// calculate offset for pagination through results
	offset := (perPage * (page - 1))

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableUser).
		Raw(dml.ListLiteUsers, perPage, offset).
		Scan(u).Error

	// variable we want to return
	users := []*library.User{}
	// iterate through all query results
	for _, user := range *u {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := user

		// convert query result to library type
		users = append(users, tmp.ToLibrary())
	}

	return users, err
}
