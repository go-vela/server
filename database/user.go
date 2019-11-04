// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// GetUser gets a user by unique ID from the database.
func (c *client) GetUser(id int64) (*library.User, error) {
	logrus.Tracef("Getting user %d from the database", id)

	// variable to store query results
	u := new(database.User)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableUser).
		Raw(c.DML.UserService.Select["user"], id).
		Scan(u).Error

	return u.ToLibrary(), err
}

// GetUserName gets a user by name from the database.
func (c *client) GetUserName(name string) (*library.User, error) {
	logrus.Tracef("Getting user %s from the database", name)

	// variable to store query results
	u := new(database.User)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableUser).
		Raw(c.DML.UserService.Select["name"], name).
		Scan(u).Error

	return u.ToLibrary(), err
}

// GetUserList gets a list of all users from the database.
func (c *client) GetUserList() ([]*library.User, error) {
	logrus.Trace("Listing users from the database")

	// variable to store query results
	u := new([]database.User)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableUser).
		Raw(c.DML.UserService.List["all"]).
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

// GetUserCount gets a list of all users from the database.
func (c *client) GetUserCount() (int64, error) {
	logrus.Trace("Counting users in the database")

	// variable to store query results
	var u []int64

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableUser).
		Raw(c.DML.UserService.Select["count"]).
		Pluck("count", &u).Error

	return u[0], err
}

// GetUserLiteList gets a lite list of all users from the database.
func (c *client) GetUserLiteList(page, perPage int) ([]*library.User, error) {
	logrus.Trace("Listing lite users from the database")

	// variable to store query results
	u := new([]database.User)
	// calculate offset for pagination through results
	offset := (perPage * (page - 1))

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableUser).
		Raw(c.DML.UserService.List["lite"], perPage, offset).
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

// CreateUser creates a new user in the database.
func (c *client) CreateUser(u *library.User) error {
	logrus.Tracef("Creating user %s from the database", u.GetName())

	// cast to database type
	user := database.UserFromLibrary(u)

	// validate the necessary fields are populated
	err := user.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Database.
		Table(constants.TableUser).
		Create(user).Error
}

// UpdateUser updates a user in the database.
func (c *client) UpdateUser(u *library.User) error {
	logrus.Tracef("Updating user %s from the database", u.GetName())

	// cast to database type
	user := database.UserFromLibrary(u)

	// validate the necessary fields are populated
	err := user.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Database.
		Table(constants.TableUser).
		Where("id = ?", u.GetID()).
		Update(user).Error
}

// DeleteUser deletes a user by unique ID from the database.
func (c *client) DeleteUser(id int64) error {
	logrus.Tracef("Deleting user %d from the database", id)

	// send query to the database
	return c.Database.
		Table(constants.TableUser).
		Raw(c.DML.UserService.Delete, id).Error
}
