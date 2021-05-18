// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"log"
	"reflect"
	"testing"

	"github.com/go-vela/server/database/sqlite/ddl"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

func init() {
	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		log.Fatalf("unable to create new sqlite test database: %v", err)
	}

	// create the user table
	err = _database.Sqlite.Exec(ddl.CreateUserTable).Error
	if err != nil {
		log.Fatalf("unable to create %s table: %v", constants.TableUser, err)
	}
}

func TestSqlite_Client_GetUser(t *testing.T) {
	// setup types
	_user := testUser()
	_user.SetID(1)
	_user.SetName("foo")
	_user.SetToken("bar")
	_user.SetHash("baz")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    *library.User
	}{
		{
			failure: false,
			want:    _user,
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the users table
		defer _database.Sqlite.Exec("delete from users;")

		// create the user in the database
		err := _database.CreateUser(test.want)
		if err != nil {
			t.Errorf("unable to create test user: %v", err)
		}

		got, err := _database.GetUser(1)

		if test.failure {
			if err == nil {
				t.Errorf("GetUser should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetUser returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetUser is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_CreateUser(t *testing.T) {
	// setup types
	_user := testUser()
	_user.SetID(1)
	_user.SetName("foo")
	_user.SetToken("bar")
	_user.SetHash("baz")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
	}{
		{
			failure: false,
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the users table
		defer _database.Sqlite.Exec("delete from users;")

		err := _database.CreateUser(_user)

		if test.failure {
			if err == nil {
				t.Errorf("CreateUser should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("CreateUser returned err: %v", err)
		}
	}
}

func TestSqlite_Client_UpdateUser(t *testing.T) {
	// setup types
	_user := testUser()
	_user.SetID(1)
	_user.SetName("foo")
	_user.SetToken("bar")
	_user.SetHash("baz")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
	}{
		{
			failure: false,
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the users table
		defer _database.Sqlite.Exec("delete from users;")

		// create the user in the database
		err := _database.CreateUser(_user)
		if err != nil {
			t.Errorf("unable to create test user: %v", err)
		}

		err = _database.UpdateUser(_user)

		if test.failure {
			if err == nil {
				t.Errorf("UpdateUser should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("UpdateUser returned err: %v", err)
		}
	}
}

func TestSqlite_Client_DeleteUser(t *testing.T) {
	// setup types
	_user := testUser()
	_user.SetID(1)
	_user.SetName("foo")
	_user.SetToken("bar")
	_user.SetHash("baz")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
	}{
		{
			failure: false,
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the users table
		defer _database.Sqlite.Exec("delete from users;")

		// create the user in the database
		err := _database.CreateUser(_user)
		if err != nil {
			t.Errorf("unable to create test user: %v", err)
		}

		err = _database.DeleteUser(1)

		if test.failure {
			if err == nil {
				t.Errorf("DeleteUser should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("DeleteUser returned err: %v", err)
		}
	}
}

// testUser is a test helper function to create a
// library User type with all fields set to their
// zero values.
func testUser() *library.User {
	i64 := int64(0)
	str := ""
	b := false
	var arr []string

	return &library.User{
		ID:           &i64,
		Name:         &str,
		RefreshToken: &str,
		Token:        &str,
		Hash:         &str,
		Favorites:    &arr,
		Active:       &b,
		Admin:        &b,
	}
}
