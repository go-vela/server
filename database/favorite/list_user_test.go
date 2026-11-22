// SPDX-License-Identifier: Apache-2.0

package favorite

import (
	"context"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestFavorite_Engine_ListUserFavorites(t *testing.T) {
	// setup types
	_owner := testutils.APIUser().Crop()
	_owner.SetID(1)
	_owner.SetName("foo")
	_owner.SetToken("bar")

	_repoOne := testutils.APIRepo()
	_repoOne.SetID(1)
	_repoOne.SetOwner(_owner)
	_repoOne.SetHash("baz")
	_repoOne.SetOrg("foo")
	_repoOne.SetName("bar")
	_repoOne.SetFullName("foo/bar")
	_repoOne.SetVisibility("public")
	_repoOne.SetPipelineType("yaml")
	_repoOne.SetTopics([]string{})
	_repoOne.SetAllowEvents(api.NewEventsFromMask(1))

	_repoTwo := testutils.APIRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetOwner(_owner)
	_repoTwo.SetHash("baz")
	_repoTwo.SetOrg("bar")
	_repoTwo.SetName("foo")
	_repoTwo.SetFullName("bar/foo")
	_repoTwo.SetVisibility("public")
	_repoTwo.SetPipelineType("yaml")
	_repoTwo.SetTopics([]string{})
	_repoTwo.SetAllowEvents(api.NewEventsFromMask(1))

	_favoriteOne := testutils.APIFavorite()
	_favoriteOne.SetRepo("foo/bar")
	_favoriteOne.SetPosition(1)

	_favoriteTwo := testutils.APIFavorite()
	_favoriteTwo.SetRepo("bar/foo")
	_favoriteTwo.SetPosition(2)

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.FavoriteFromAPI(_favoriteOne), *types.FavoriteFromAPI(_favoriteTwo)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT r.full_name as repo_name, f.position FROM favorites f JOIN repos r ON f.repo_id = r.id WHERE f.user_id = $1 ORDER BY f.position ASC`).WillReturnRows(_rows)

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

	err = _sqlite.CreateFavorite(context.TODO(), _owner, _favoriteOne)
	if err != nil {
		t.Errorf("unable to create test favorite for sqlite: %v", err)
	}

	err = _sqlite.CreateFavorite(context.TODO(), _owner, _favoriteTwo)
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
			want:     []*api.Favorite{_favoriteOne, _favoriteTwo},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.Favorite{_favoriteOne, _favoriteTwo},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListUserFavorites(context.TODO(), _owner)

			if test.failure {
				if err == nil {
					t.Errorf("ListRepos for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListRepos for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListRepos for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
