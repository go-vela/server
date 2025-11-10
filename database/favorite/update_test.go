// SPDX-License-Identifier: Apache-2.0

package favorite

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestFavorite_Engine_UpdateFavorites(t *testing.T) {
	// setup types
	_user := testutils.APIUser()
	_user.SetID(1)

	_repoOne := testutils.APIRepo()
	_repoOne.SetID(1)
	_repoOne.SetFullName("foo/bar")

	_repoTwo := testutils.APIRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetFullName("bar/foo")

	_favoriteOne := testutils.APIFavorite()
	_favoriteOne.SetRepo("foo/bar")
	_favoriteOne.SetPosition(1)

	_favoriteTwo := testutils.APIFavorite()
	_favoriteTwo.SetRepo("bar/foo")
	_favoriteTwo.SetPosition(2)

	_updatedFavoriteOne := testutils.APIFavorite()
	_updatedFavoriteOne.SetRepo("foo/bar")
	_updatedFavoriteOne.SetPosition(2)

	_updatedFavoriteTwo := testutils.APIFavorite()
	_updatedFavoriteTwo.SetRepo("bar/foo")
	_updatedFavoriteTwo.SetPosition(1)

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`WITH input(repo_name, position) AS ( VALUES ($1, $2), ($3, $4) ) INSERT INTO favorites (user_id, repo_id, position) SELECT $5, r.id, input.position FROM input JOIN repos r ON r.full_name = input.repo_name ON CONFLICT (user_id, repo_id) DO UPDATE SET position = excluded.position;`).
		WithArgs(1, "foo/bar", 2, "bar/foo", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_rows := testutils.CreateMockRows([]any{*types.FavoriteFromAPI(_updatedFavoriteTwo), *types.FavoriteFromAPI(_updatedFavoriteOne)})

	_mock.ExpectQuery(`SELECT r.full_name as repo_name, f.position FROM favorites f JOIN repos r ON f.repo_id = r.id WHERE f.user_id = $1 ORDER BY f.position ASC`).WillReturnRows(_rows)

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.client.AutoMigrate(&types.User{})
	if err != nil {
		t.Errorf("unable to create user table for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableUser).Create(types.UserFromAPI(_user)).Error
	if err != nil {
		t.Errorf("unable to create test user for sqlite: %v", err)
	}

	err = _sqlite.client.AutoMigrate(&types.Repo{})
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

	err = _sqlite.CreateFavorite(context.TODO(), _user, _favoriteOne)
	if err != nil {
		t.Errorf("unable to create test favorite for sqlite: %v", err)
	}

	err = _sqlite.CreateFavorite(context.TODO(), _user, _favoriteTwo)
	if err != nil {
		t.Errorf("unable to create test favorite for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     []*api.Favorite
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*api.Favorite{_updatedFavoriteTwo, _updatedFavoriteOne},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.Favorite{_updatedFavoriteTwo, _updatedFavoriteOne},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.database.UpdateFavorites(context.TODO(), _user, []*api.Favorite{_updatedFavoriteOne, _updatedFavoriteTwo})

			if test.failure {
				if err == nil {
					t.Errorf("UpdateFavorites for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateFavorites for %s returned err: %v", test.name, err)
			}

			list, err := test.database.ListUserFavorites(context.TODO(), _user)
			if err != nil {
				t.Errorf("ListUserFavorites for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(list, test.want) {
				t.Errorf("ListUserFavorites for %s returned %v, want %v", test.name, list, test.want)
			}
		})
	}
}
