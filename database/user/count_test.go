// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestUser_Engine_CountUsers(t *testing.T) {
	// setup types
	_userOne := testutils.APIUser()
	_userOne.SetID(1)
	_userOne.SetName("foo")
	_userOne.SetToken("bar")

	_userTwo := testutils.APIUser()
	_userTwo.SetID(2)
	_userTwo.SetName("baz")
	_userTwo.SetToken("bar")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "users"`).WillReturnRows(_rows)

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
		database *Engine
		want     int64
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     2,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     2,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.CountUsers(context.TODO())

			if test.failure {
				if err == nil {
					t.Errorf("CountUsers for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CountUsers for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("CountUsers for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
