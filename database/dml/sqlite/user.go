// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

const (
	// ListUsers represents a query to
	// list all users in the database.
	ListUsers = `
SELECT *
FROM users;
`

	// ListLiteUsers represents a query to
	// list all lite users in the database.
	ListLiteUsers = `
SELECT id, name
FROM users
ORDER BY id DESC
LIMIT ?
OFFSET ?;
`

	// SelectUser represents a query to select
	// a user for an id in the database.
	SelectUser = `
SELECT *
FROM users
WHERE id = ?
LIMIT 1;
`

	// SelectUserName represents a query to select
	// a user for a name in the database.
	SelectUserName = `
SELECT *
FROM users
WHERE name = ?
LIMIT 1;
`

	// SelectUsersCount represents a query to select
	// the count of users in the database.
	SelectUsersCount = `
SELECT count(*) as count
FROM users;
`

	// SelectRefreshToken represents a query to select
	// a user for a refresh_token in the database.
	//
	// nolint: gosec // ignore false positive
	SelectRefreshToken = `
SELECT *
FROM users
WHERE refresh_token = ?
LIMIT 1;
`

	// DeleteUser represents a query to
	// remove a user from the database.
	DeleteUser = `
DELETE
FROM users
WHERE id = ?;
`
)

// createUserService is a helper function to return
// a service for interacting with the users table.
func createUserService() *Service {
	return &Service{
		List: map[string]string{
			"all":  ListUsers,
			"lite": ListLiteUsers,
		},
		Select: map[string]string{
			"user":         SelectUser,
			"name":         SelectUserName,
			"count":        SelectUsersCount,
			"refreshToken": SelectRefreshToken,
		},
		Delete: DeleteUser,
	}
}
