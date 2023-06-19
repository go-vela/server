// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package user

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestUser_Engine_GetUserForName(t *testing.T) {
	// setup types
	_user := testUser()
	_user.SetID(1)
	_user.SetName("foo")
	_user.SetToken("bar")
	_user.SetHash("baz")
	_user.SetFavorites([]string{})

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "name", "refresh_token", "token", "hash", "favorites", "active", "admin"}).
		AddRow(1, "foo", "", "bar", "baz", "{}", false, false)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "users" WHERE name = $1 LIMIT 1`).WithArgs("foo").WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateUser(_user)
	if err != nil {
		t.Errorf("unable to create test user for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     *library.User
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _user,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _user,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetUserForName("foo")

			if test.failure {
				if err == nil {
					t.Errorf("GetUserForName for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetUserForName for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetUserForName for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
