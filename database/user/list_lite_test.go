// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
)

func TestUser_Engine_ListLiteUsers(t *testing.T) {
	// setup types
	_userOne := testutils.APIUser()
	_userOne.SetID(1)
	_userOne.SetName("foo")
	_userOne.SetToken("bar")
	_userOne.SetFavorites([]string{})
	_userOne.SetDashboards([]string{})

	_userTwo := testutils.APIUser()
	_userTwo.SetID(2)
	_userTwo.SetName("baz")
	_userTwo.SetToken("bar")
	_userTwo.SetFavorites([]string{})
	_userTwo.SetDashboards([]string{})

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "name"}).
		AddRow(1, "foo").
		AddRow(2, "baz")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT "id","name" FROM "users" LIMIT $1`).WithArgs(10).WillReturnRows(_rows)

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

	// empty fields not returned by query
	_userOne.RefreshToken = new(string)
	_userOne.Token = new(string)
	_userOne.Favorites = new([]string)
	_userOne.Dashboards = new([]string)

	_userTwo.RefreshToken = new(string)
	_userTwo.Token = new(string)
	_userTwo.Favorites = new([]string)
	_userTwo.Dashboards = new([]string)

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     []*api.User
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*api.User{_userOne, _userTwo},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.User{_userTwo, _userOne},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListLiteUsers(context.TODO(), 1, 10)

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
