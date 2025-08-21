// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestBuild_Engine_ListBuildsForOrg(t *testing.T) {
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
	_repoOne.SetAllowEvents(api.NewEventsFromMask(1))
	_repoOne.SetTopics([]string{})

	_repoTwo := testutils.APIRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetOwner(_owner)
	_repoTwo.SetHash("bar")
	_repoTwo.SetOrg("foo")
	_repoTwo.SetName("baz")
	_repoTwo.SetFullName("foo/baz")
	_repoTwo.SetVisibility("public")
	_repoTwo.SetPipelineType("yaml")
	_repoTwo.SetAllowEvents(api.NewEventsFromMask(1))
	_repoTwo.SetTopics([]string{})

	_buildOne := testutils.APIBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepo(_repoOne)
	_buildOne.SetNumber(1)
	_buildOne.SetDeployPayload(nil)
	_buildOne.SetEvent("push")

	_buildTwo := testutils.APIBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepo(_repoTwo)
	_buildTwo.SetNumber(2)
	_buildTwo.SetDeployPayload(nil)
	_buildTwo.SetEvent("push")

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected query without filters result in mock
	_rows := testutils.CreateMockRows([]any{*types.BuildFromAPI(_buildOne), *types.BuildFromAPI(_buildTwo)})

	_repoRows := testutils.CreateMockRows([]any{*types.RepoFromAPI(_repoOne), *types.RepoFromAPI(_repoTwo)})

	_userRows := testutils.CreateMockRows([]any{*types.UserFromAPI(_owner)})

	// ensure the mock expects the query without filters
	_mock.ExpectQuery(`SELECT builds.* FROM "builds" JOIN repos ON builds.repo_id = repos.id WHERE repos.org = $1 ORDER BY created DESC,id LIMIT $2`).WithArgs("foo", 10).WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE "repos"."id" IN ($1,$2)`).WithArgs(1, 2).WillReturnRows(_repoRows)
	_mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."id" = $1`).WithArgs(1).WillReturnRows(_userRows)

	// create expected query with event filter result in mock
	_rows = testutils.CreateMockRows([]any{*types.BuildFromAPI(_buildOne), *types.BuildFromAPI(_buildTwo)})
	_repoRows = testutils.CreateMockRows([]any{*types.RepoFromAPI(_repoOne), *types.RepoFromAPI(_repoTwo)})
	_userRows = testutils.CreateMockRows([]any{*types.UserFromAPI(_owner)})

	// ensure the mock expects the query with event filter
	_mock.ExpectQuery(`SELECT builds.* FROM "builds" JOIN repos ON builds.repo_id = repos.id WHERE repos.org = $1 AND builds.event = $2 ORDER BY created DESC,id LIMIT $3`).WithArgs("foo", "push", 10).WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE "repos"."id" IN ($1,$2)`).WithArgs(1, 2).WillReturnRows(_repoRows)
	_mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."id" = $1`).WithArgs(1).WillReturnRows(_userRows)

	// create expected query with visibility filter result in mock
	_rows = testutils.CreateMockRows([]any{*types.BuildFromAPI(_buildOne), *types.BuildFromAPI(_buildTwo)})
	_repoRows = testutils.CreateMockRows([]any{*types.RepoFromAPI(_repoOne), *types.RepoFromAPI(_repoTwo)})
	_userRows = testutils.CreateMockRows([]any{*types.UserFromAPI(_owner)})

	// ensure the mock expects the query with visibility filter
	_mock.ExpectQuery(`SELECT builds.* FROM "builds" JOIN repos ON builds.repo_id = repos.id WHERE repos.org = $1 AND repos.visibility = $2 ORDER BY created DESC,id LIMIT $3`).WithArgs("foo", "public", 10).WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE "repos"."id" IN ($1,$2)`).WithArgs(1, 2).WillReturnRows(_repoRows)
	_mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."id" = $1`).WithArgs(1).WillReturnRows(_userRows)

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateBuild(context.TODO(), _buildOne)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	_, err = _sqlite.CreateBuild(context.TODO(), _buildTwo)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	err = _sqlite.client.AutoMigrate(&types.Repo{})
	if err != nil {
		t.Errorf("unable to create repo table for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableRepo).Create(types.RepoFromAPI(_repoOne)).Error
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableRepo).Create(types.RepoFromAPI(_repoTwo)).Error
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	err = _sqlite.client.AutoMigrate(&types.User{})
	if err != nil {
		t.Errorf("unable to create build table for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableUser).Create(types.UserFromAPI(_owner)).Error
	if err != nil {
		t.Errorf("unable to create test user for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure      bool
		name         string
		database     *Engine
		repoFilters  map[string]any
		buildFilters map[string]any
		want         []*api.Build
	}{
		{
			failure:  false,
			name:     "postgres without filters",
			database: _postgres,
			want:     []*api.Build{_buildOne, _buildTwo},
		},
		{
			failure:  false,
			name:     "postgres with event filter",
			database: _postgres,
			buildFilters: map[string]any{
				"event": "push",
			},
			want: []*api.Build{_buildOne, _buildTwo},
		},
		{
			failure:  false,
			name:     "postgres with visibility filter",
			database: _postgres,
			repoFilters: map[string]any{
				"visibility": "public",
			},
			want: []*api.Build{_buildOne, _buildTwo},
		},
		{
			failure:  false,
			name:     "sqlite3 without filters",
			database: _sqlite,
			want:     []*api.Build{_buildOne, _buildTwo},
		},
		{
			failure:  false,
			name:     "sqlite3 with event filter",
			database: _sqlite,
			buildFilters: map[string]any{
				"event": "push",
			},
			want: []*api.Build{_buildOne, _buildTwo},
		},
		{
			failure:  false,
			name:     "sqlite3 with visibility filter",
			database: _sqlite,
			repoFilters: map[string]any{
				"visibility": "public",
			},
			want: []*api.Build{_buildOne, _buildTwo},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListBuildsForOrg(context.TODO(), "foo", test.repoFilters, test.buildFilters, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListBuildsForOrg for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListBuildsForOrg for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("ListBuildsForOrg for %s mismatch (-want +got):\n%s", test.name, diff)
			}
		})
	}
}
