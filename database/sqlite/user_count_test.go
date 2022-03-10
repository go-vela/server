// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"log"
	"reflect"
	"testing"

	"github.com/go-vela/server/database/sqlite/ddl"
	"github.com/go-vela/types/constants"
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

func TestSqlite_Client_GetUserCount(t *testing.T) {
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
		want    int64
	}{
		{
			failure: false,
			want:    2,
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the users table
		defer _database.Sqlite.Exec("delete from users;")

		// create the users in the database
		err := _database.CreateUser(_userOne)
		if err != nil {
			t.Errorf("unable to create test user: %v", err)
		}

		err = _database.CreateUser(_userTwo)
		if err != nil {
			t.Errorf("unable to create test user: %v", err)
		}

		got, err := _database.GetUserCount()

		if test.failure {
			if err == nil {
				t.Errorf("GetUserCount should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetUserCount returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetUserCount is %v, want %v", got, test.want)
		}
	}
}
