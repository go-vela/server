// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package user

import (
	"github.com/go-vela/types/library"
)

// UserService represents the Vela interface for user
// functions with the supported Database backends.
//
// nolint: revive // ignore name stutter
type UserService interface {
	// User Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateUserIndexes defines a function that creates the indexes for the users table.
	CreateUserIndexes() error
	// CreateUserTable defines a function that creates the users table.
	CreateUserTable(string) error

	// User Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CountUsers defines a function that gets the count of all users.
	CountUsers() (int64, error)
	// CreateUser defines a function that creates a new user.
	CreateUser(*library.User) error
	// DeleteUser defines a function that deletes an existing user.
	DeleteUser(*library.User) error
	// GetUser defines a function that gets a user by ID.
	GetUser(int64) (*library.User, error)
	// GetUserForName defines a function that gets a user by name.
	GetUserForName(string) (*library.User, error)
	// ListUsers defines a function that gets a list of all users.
	ListUsers() ([]*library.User, error)
	// ListLiteUsers defines a function that gets a lite list of users.
	ListLiteUsers(int, int) ([]*library.User, int64, error)
	// UpdateUser defines a function that updates an existing user.
	UpdateUser(*library.User) error
}
