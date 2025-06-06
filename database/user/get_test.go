// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestUser_Engine_GetUser(t *testing.T) {
	// setup types
	_user := testutils.APIUser()
	_user.SetID(1)
	_user.SetName("foo")
	_user.SetToken("bar")

	_user.SetFavorites([]string{})
	_user.SetDashboards([]string{})

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.UserFromAPI(_user)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "users" WHERE id = $1 LIMIT $2`).WithArgs(1, 1).WillReturnRows(_rows)

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
		database *Engine
		want     *api.User
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
			got, err := test.database.GetUser(context.TODO(), 1)

			if test.failure {
				if err == nil {
					t.Errorf("GetUser for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetUser for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetUser for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
