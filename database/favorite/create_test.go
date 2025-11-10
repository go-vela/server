// SPDX-License-Identifier: Apache-2.0

package favorite

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestFavorite_Engine_CreateFavorite(t *testing.T) {
	// setup types
	_user := testutils.APIUser()
	_user.SetID(1)

	_repo := testutils.APIRepo()
	_repo.SetID(1)
	_repo.SetFullName("foo/bar")

	_favorite := testutils.APIFavorite()
	_favorite.SetRepo("foo/bar")
	_favorite.SetPosition(1)

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`INSERT INTO favorites
(user_id, repo_id, position)
SELECT $1, id, $2 FROM repos WHERE full_name = $3;`).
		WithArgs(1, 1, "foo/bar").WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)

	err := _sqlite.client.AutoMigrate(&types.Repo{})
	if err != nil {
		t.Errorf("unable to create build table for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableRepo).Create(types.RepoFromAPI(_repo)).Error
	if err != nil {
		t.Errorf("unable to create test user for sqlite: %v", err)
	}

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
			err := test.database.CreateFavorite(context.TODO(), _user, _favorite)

			if test.failure {
				if err == nil {
					t.Errorf("CreateUser for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateUser for %s returned err: %v", test.name, err)
			}
		})
	}
}
