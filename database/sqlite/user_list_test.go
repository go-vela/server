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

func TestSqlite_Client_GetUserList(t *testing.T) {
	// setup types
	_userOne := testUser()
	_userOne.SetID(1)
	_userOne.SetName("foo")
	_userOne.SetToken("bar")
	_userOne.SetHash("baz")

	_userTwo := testUser()
	_userTwo.SetID(2)
	_userTwo.SetName("bar")
	_userTwo.SetToken("foo")
	_userTwo.SetHash("baz")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.User
	}{
		{
			failure: false,
			want:    []*library.User{_userOne, _userTwo},
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the users table
		defer _database.Sqlite.Exec("delete from users;")

		for _, user := range test.want {
			// create the user in the database
			err := _database.CreateUser(user)
			if err != nil {
				t.Errorf("unable to create test user: %v", err)
			}
		}

		got, err := _database.GetUserList()

		if test.failure {
			if err == nil {
				t.Errorf("GetUserList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetUserList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetUserList is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetUserLiteList(t *testing.T) {
	// setup types
	_userOne := testUser()
	_userOne.SetID(1)
	_userOne.SetName("foo")

	_userTwo := testUser()
	_userTwo.SetID(2)
	_userTwo.SetName("bar")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.User
	}{
		{
			failure: false,
			want:    []*library.User{_userTwo, _userOne},
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the users table
		defer _database.Sqlite.Exec("delete from users;")

		for _, user := range test.want {
			// set the required fields for the user
			user.SetToken("baz")
			user.SetHash("foob")

			// create the user in the database
			err := _database.CreateUser(user)
			if err != nil {
				t.Errorf("unable to create test user: %v", err)
			}

			// clear the required fields for the user
			// so we get back the expected data
			user.SetToken("")
			user.SetHash("")
		}

		got, err := _database.GetUserLiteList(1, 10)

		if test.failure {
			if err == nil {
				t.Errorf("GetUserLiteList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetUserLiteList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetUserLiteList is %v, want %v", got, test.want)
		}
	}
}
