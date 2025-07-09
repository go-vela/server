// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestUser_Engine_CreateUser(t *testing.T) {
	// setup types
	_user := testutils.APIUser()
	_user.SetID(1)
	_user.SetName("foo")
	_user.SetToken("bar")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "users"
("name","refresh_token","token","favorites","active","admin","dashboards","id")
VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`).
		WithArgs("foo", AnyArgument{}, AnyArgument{}, nil, false, false, AnyArgument{}, 1).
		WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
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
			got, err := test.database.CreateUser(context.TODO(), _user)

			if test.failure {
				if err == nil {
					t.Errorf("CreateUser for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateUser for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _user) {
				t.Errorf("CreateUser for %s returned %s, want %s", test.name, got, _user)
			}
		})
	}
}
