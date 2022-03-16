// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/postgres/dml"
	"github.com/go-vela/types/library"

	"gorm.io/gorm"
)

func TestPostgres_Client_GetUserList(t *testing.T) {
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
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.ListUsers).Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "name", "refresh_token", "token", "hash", "favorites", "active", "admin"},
	).AddRow(1, "foo", "", "bar", "baz", "{}", false, false).
		AddRow(2, "bar", "", "foo", "baz", "{}", false, false)

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

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

func TestPostgres_Client_GetUserLiteList(t *testing.T) {
	// setup types
	_userOne := testUser()
	_userOne.SetID(1)
	_userOne.SetName("foo")
	_userOne.SetFavorites(nil)

	_userTwo := testUser()
	_userTwo.SetID(2)
	_userTwo.SetName("bar")
	_userTwo.SetFavorites(nil)

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.ListLiteUsers, 1, 10).Statement

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "foo").AddRow(2, "bar")

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

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
