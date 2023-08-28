// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package user

import (
	"context"

	"github.com/go-vela/types/library"
)

// UserInterface represents the Vela interface for user
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type UserInterface interface {
	// User Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateUserIndexes defines a function that creates the indexes for the users table.
	CreateUserIndexes(context.Context) error
	// CreateUserTable defines a function that creates the users table.
	CreateUserTable(context.Context, string) error

	// User Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CountUsers defines a function that gets the count of all users.
	CountUsers(context.Context) (int64, error)
	// CreateUser defines a function that creates a new user.
	CreateUser(context.Context, *library.User) error
	// DeleteUser defines a function that deletes an existing user.
	DeleteUser(context.Context, *library.User) error
	// GetUser defines a function that gets a user by ID.
	GetUser(context.Context, int64) (*library.User, error)
	// GetUserForName defines a function that gets a user by name.
	GetUserForName(context.Context, string) (*library.User, error)
	// ListUsers defines a function that gets a list of all users.
	ListUsers(context.Context) ([]*library.User, error)
	// ListLiteUsers defines a function that gets a lite list of users.
	ListLiteUsers(context.Context, int, int) ([]*library.User, int64, error)
	// UpdateUser defines a function that updates an existing user.
	UpdateUser(context.Context, *library.User) error
}
