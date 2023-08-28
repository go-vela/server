// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package user

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestUser_Engine_ListLiteUsers(t *testing.T) {
	// setup types
	_userOne := testUser()
	_userOne.SetID(1)
	_userOne.SetName("foo")
	_userOne.SetToken("bar")
	_userOne.SetHash("baz")
	_userOne.SetFavorites([]string{})

	_userTwo := testUser()
	_userTwo.SetID(2)
	_userTwo.SetName("baz")
	_userTwo.SetToken("bar")
	_userTwo.SetHash("foo")
	_userTwo.SetFavorites([]string{})

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "users"`).WillReturnRows(_rows)

	// create expected result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "name"}).
		AddRow(1, "foo").
		AddRow(2, "baz")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT "id","name" FROM "users" LIMIT 10`).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateUser(context.TODO(), _userOne)
	if err != nil {
		t.Errorf("unable to create test user for sqlite: %v", err)
	}

	err = _sqlite.CreateUser(context.TODO(), _userTwo)
	if err != nil {
		t.Errorf("unable to create test user for sqlite: %v", err)
	}

	// empty fields not returned by query
	_userOne.RefreshToken = new(string)
	_userOne.Token = new(string)
	_userOne.Hash = new(string)
	_userOne.Favorites = new([]string)

	_userTwo.RefreshToken = new(string)
	_userTwo.Token = new(string)
	_userTwo.Hash = new(string)
	_userTwo.Favorites = new([]string)

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     []*library.User
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*library.User{_userOne, _userTwo},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*library.User{_userTwo, _userOne},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _, err := test.database.ListLiteUsers(context.TODO(), 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListLiteUsers for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListLiteUsers for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListLiteUsers for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
