// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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

func TestPostgres_Client_GetUser(t *testing.T) {
	// setup types
	_user := testUser()
	_user.SetID(1)
	_user.SetName("foo")
	_user.SetToken("bar")
	_user.SetHash("baz")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.SelectUser, 1).Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "name", "refresh_token", "token", "hash", "favorites", "active", "admin"},
	).AddRow(1, "foo", "", "bar", "baz", "{}", false, false)

	// ensure the mock expects the query for test case 1
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)
	// ensure the mock expects the error for test case 2
	_mock.ExpectQuery(_query.SQL.String()).WillReturnError(gorm.ErrRecordNotFound)

	// setup tests
	tests := []struct {
		failure bool
		want    *library.User
	}{
		{
			failure: false,
			want:    _user,
		},
		{
			failure: true,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
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

func TestPostgres_Client_CreateUser(t *testing.T) {
	// setup types
	_user := testUser()
	_user.SetID(1)
	_user.SetName("foo")
	_user.SetToken("bar")
	_user.SetHash("baz")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`INSERT INTO "users" ("name","refresh_token","token","hash","favorites","active","admin","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`).
		WithArgs("foo", AnyArgument{}, AnyArgument{}, AnyArgument{}, "{}", false, false, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

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

func TestPostgres_Client_UpdateUser(t *testing.T) {
	// setup types
	_user := testUser()
	_user.SetID(1)
	_user.SetName("foo")
	_user.SetToken("bar")
	_user.SetHash("baz")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "users" SET "name"=$1,"refresh_token"=$2,"token"=$3,"hash"=$4,"favorites"=$5,"active"=$6,"admin"=$7 WHERE "id" = $8`).
		WithArgs("foo", AnyArgument{}, AnyArgument{}, AnyArgument{}, "{}", false, false, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

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
		err := _database.UpdateUser(_user)

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

func TestPostgres_Client_DeleteUser(t *testing.T) {
	// setup types

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Exec(dml.DeleteUser, 1).Statement

	// ensure the mock expects the query
	_mock.ExpectExec(_query.SQL.String()).WillReturnResult(sqlmock.NewResult(1, 1))

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
		err := _database.DeleteUser(1)

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
	arr := []string{}
	b := false

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
