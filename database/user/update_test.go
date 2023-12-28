// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUser_Engine_UpdateUser(t *testing.T) {
	// setup types
	_user := testUser()
	_user.SetID(1)
	_user.SetName("foo")
	_user.SetToken("bar")
	_user.SetHash("baz")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "users"
SET "name"=$1,"refresh_token"=$2,"token"=$3,"hash"=$4,"favorites"=$5,"active"=$6,"admin"=$7
WHERE "id" = $8`).
		WithArgs("foo", AnyArgument{}, AnyArgument{}, AnyArgument{}, nil, false, false, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateUser(context.TODO(), _user)
	if err != nil {
		t.Errorf("unable to create test user for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.UpdateUser(context.TODO(), _user)

			if test.failure {
				if err == nil {
					t.Errorf("UpdateUser for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateUser for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _user) {
				t.Errorf("UpdateUser for %s returned %s, want %s", test.name, got, _user)
			}
		})
	}
}
