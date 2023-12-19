// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestUser_Engine_ListUsers(t *testing.T) {
	// setup types
	_userOne := testUser()
	_userOne.SetID(1)
	_userOne.SetName("foo")
	_userOne.SetToken("bar")
	_userOne.SetHash("baz")
	_userOne.SetFavorites([]string{})
	_userOne.SetDashboards([]string{})

	_userTwo := testUser()
	_userTwo.SetID(2)
	_userTwo.SetName("baz")
	_userTwo.SetToken("bar")
	_userTwo.SetHash("foo")
	_userTwo.SetFavorites([]string{})
	_userTwo.SetDashboards([]string{})

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "users"`).WillReturnRows(_rows)

	// create expected result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "name", "refresh_token", "token", "hash", "favorites", "active", "admin", "dashboards"}).
		AddRow(1, "foo", "", "bar", "baz", "{}", false, false, "{}").
		AddRow(2, "baz", "", "bar", "foo", "{}", false, false, "{}")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "users"`).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateUser(context.TODO(), _userOne)
	if err != nil {
		t.Errorf("unable to create test user for sqlite: %v", err)
	}

	_, err = _sqlite.CreateUser(context.TODO(), _userTwo)
	if err != nil {
		t.Errorf("unable to create test user for sqlite: %v", err)
	}

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
			want:     []*library.User{_userOne, _userTwo},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListUsers(context.TODO())

			if test.failure {
				if err == nil {
					t.Errorf("ListUsers for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListUsers for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListUsers for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
