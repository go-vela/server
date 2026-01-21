// SPDX-License-Identifier: Apache-2.0

package favorite

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestFavorite_Engine_UpdatePosition(t *testing.T) {
	// setup types
	_favoriteOne := testutils.APIFavorite()
	_favoriteOne.SetRepo("foo/bar")
	_favoriteOne.SetPosition(1)

	_favoriteTwo := testutils.APIFavorite()
	_favoriteTwo.SetRepo("baz/qux")
	_favoriteTwo.SetPosition(2)

	_favoriteThree := testutils.APIFavorite()
	_favoriteThree.SetRepo("quux/corge")
	_favoriteThree.SetPosition(3)

	_user := testutils.APIUser()
	_user.SetID(1)

	_repoOne := testutils.APIRepo()
	_repoOne.SetID(1)
	_repoOne.SetFullName("foo/bar")

	_repoTwo := testutils.APIRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetFullName("baz/qux")

	_repoThree := testutils.APIRepo()
	_repoThree.SetID(3)
	_repoThree.SetFullName("quux/corge")

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// move down mocks
	_mock.ExpectBegin()

	_mock.ExpectQuery(`SELECT position FROM "favorites" WHERE user_id = $1 AND repo_id = $2`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"position"}).AddRow(1))

	_mock.ExpectQuery(`SELECT count(*) FROM "favorites" WHERE user_id = $1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	_mock.ExpectExec(`UPDATE favorites SET position = position - 1 WHERE user_id = $1 AND position <= $2 AND position > $3`).
		WithArgs(1, 2, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_mock.ExpectExec(`UPDATE "favorites" SET "position"=$1 WHERE user_id = $2 AND repo_id = $3`).
		WithArgs(2, 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_mock.ExpectCommit()

	// move up mocks
	_mock.ExpectBegin()

	_mock.ExpectQuery(`SELECT position FROM "favorites" WHERE user_id = $1 AND repo_id = $2`).
		WithArgs(1, 3).
		WillReturnRows(sqlmock.NewRows([]string{"position"}).AddRow(3))

	_mock.ExpectQuery(`SELECT count(*) FROM "favorites" WHERE user_id = $1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	_mock.ExpectExec(`UPDATE favorites SET position = position + 1 WHERE user_id = $1 AND position >= $2 AND position < $3`).
		WithArgs(1, 1, 3).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_mock.ExpectExec(`UPDATE "favorites" SET "position"=$1 WHERE user_id = $2 AND repo_id = $3`).
		WithArgs(1, 1, 3).
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

	err = _sqlite.client.Table(constants.TableRepo).Create(types.RepoFromAPI(_repoThree)).Error
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

	err = _sqlite.CreateFavorite(t.Context(), _user, _favoriteThree)
	if err != nil {
		t.Errorf("unable to create test favorite for sqlite: %v", err)
	}

	int2 := int64(2)
	int1 := int64(1)

	// setup tests
	tests := []struct {
		failure    bool
		name       string
		database   *Engine
		updateRepo *api.Repo
		updateFav  *api.Favorite
	}{
		{
			failure:    false,
			name:       "postgres move down",
			database:   _postgres,
			updateRepo: _repoOne,
			updateFav: &api.Favorite{
				Position: &int2,
			},
		},
		{
			failure:    false,
			name:       "postgres move up",
			database:   _postgres,
			updateRepo: _repoThree,
			updateFav: &api.Favorite{
				Position: &int1,
			},
		},
		{
			failure:    false,
			name:       "sqlite3 move down",
			database:   _sqlite,
			updateRepo: _repoOne,
			updateFav: &api.Favorite{
				Position: &int2,
			},
		},
		{
			failure:    false,
			name:       "sqlite3 move up",
			database:   _sqlite,
			updateRepo: _repoThree,
			updateFav: &api.Favorite{
				Position: &int1,
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.database.UpdateFavoritePosition(t.Context(), _user, test.updateRepo, test.updateFav)

			if test.failure {
				if err == nil {
					t.Errorf("UpdateFavoritePosition for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateFavoritePosition for %s returned err: %v", test.name, err)
			}
		})
	}
}
