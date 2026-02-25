// SPDX-License-Identifier: Apache-2.0

package favorite

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestFavorite_Engine_DeleteFavorite(t *testing.T) {
	// setup types
	_favoriteOne := testutils.APIFavorite()
	_favoriteOne.SetRepo("foo/bar")
	_favoriteOne.SetPosition(1)

	_favoriteTwo := testutils.APIFavorite()
	_favoriteTwo.SetRepo("baz/qux")
	_favoriteTwo.SetPosition(2)

	_user := testutils.APIUser()
	_user.SetID(1)

	_repoOne := testutils.APIRepo()
	_repoOne.SetID(1)
	_repoOne.SetFullName("foo/bar")

	_repoTwo := testutils.APIRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetFullName("baz/qux")

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	_mock.ExpectBegin()

	_mock.ExpectQuery(`SELECT position FROM "favorites" WHERE user_id = $1 AND repo_id = $2`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"position"}).AddRow(1))

	_mock.ExpectExec(`DELETE FROM favorites WHERE user_id = $1 AND repo_id = $2`).WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))

	_mock.ExpectExec(`UPDATE favorites SET position = position - 1 WHERE user_id = $1 AND position > $2`).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_mock.ExpectCommit()

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.client.AutoMigrate(&types.Repo{})
	if err != nil {
		t.Errorf("unable to create build table for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableRepo).Create(types.RepoFromAPI(_repoOne)).Error
	if err != nil {
		t.Errorf("unable to create test user for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableRepo).Create(types.RepoFromAPI(_repoTwo)).Error
	if err != nil {
		t.Errorf("unable to create test user for sqlite: %v", err)
	}

	err = _sqlite.CreateFavorite(t.Context(), _user, _favoriteOne)
	if err != nil {
		t.Errorf("unable to create test favorite for sqlite: %v", err)
	}

	err = _sqlite.CreateFavorite(t.Context(), _user, _favoriteTwo)
	if err != nil {
		t.Errorf("unable to create test favorite for sqlite: %v", err)
	}

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
			err := test.database.DeleteFavorite(t.Context(), _user, _repoOne)

			if test.failure {
				if err == nil {
					t.Errorf("DeleteFavorite for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("DeleteFavorite for %s returned err: %v", test.name, err)
			}
		})
	}
}
