// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestRepo_Engine_ListReposForOrg(t *testing.T) {
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
	_repoTwo.SetHash("bar")
	_repoTwo.SetOrg("foo")
	_repoTwo.SetName("baz")
	_repoTwo.SetFullName("foo/baz")
	_repoTwo.SetVisibility("public")
	_repoTwo.SetPipelineType("yaml")
	_repoTwo.SetTopics([]string{})
	_repoTwo.SetAllowEvents(api.NewEventsFromMask(1))

	_buildOne := new(api.Build)
	_buildOne.SetID(1)
	_buildOne.SetRepo(_repoOne)
	_buildOne.SetNumber(1)
	_buildOne.SetCreated(time.Now().UTC().Unix())

	_buildTwo := new(api.Build)
	_buildTwo.SetID(2)
	_buildTwo.SetRepo(_repoTwo)
	_buildTwo.SetNumber(1)
	_buildTwo.SetCreated(time.Now().UTC().Unix())

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected name query result in mock
	_rows := testutils.CreateMockRows([]any{*types.RepoFromAPI(_repoOne), *types.RepoFromAPI(_repoTwo)})

	_userRows := testutils.CreateMockRows([]any{*types.UserFromAPI(_owner)})

	// ensure the mock expects the name query
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE org = $1 ORDER BY name LIMIT $2`).WithArgs("foo", 10).WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."id" = $1`).WithArgs(1).WillReturnRows(_userRows)

	// create expected latest query result in mock
	_rows = testutils.CreateMockRows([]any{*types.RepoFromAPI(_repoOne), *types.RepoFromAPI(_repoTwo)})

	_userRows = testutils.CreateMockRows([]any{*types.UserFromAPI(_owner)})

	// ensure the mock expects the latest query
	_mock.ExpectQuery(`SELECT repos.* FROM "repos" LEFT JOIN (SELECT repos.id, MAX(builds.created) AS latest_build FROM "builds" INNER JOIN repos repos ON builds.repo_id = repos.id WHERE repos.org = $1 GROUP BY "repos"."id") t on repos.id = t.id ORDER BY latest_build DESC NULLS LAST LIMIT $2`).WithArgs("foo", 10).WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."id" = $1`).WithArgs(1).WillReturnRows(_userRows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateRepo(context.TODO(), _repoOne)
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	_, err = _sqlite.CreateRepo(context.TODO(), _repoTwo)
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	err = _sqlite.client.Migrator().CreateTable(&types.User{})
	if err != nil {
		t.Errorf("unable to create build table for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableUser).Create(types.UserFromAPI(_owner)).Error
	if err != nil {
		t.Errorf("unable to create test user for sqlite: %v", err)
	}

	err = _sqlite.client.Migrator().CreateTable(&types.Build{})
	if err != nil {
		t.Errorf("unable to create build table for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableBuild).Create(types.BuildFromAPI(_buildOne)).Error
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableBuild).Create(types.BuildFromAPI(_buildTwo)).Error
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		sort     string
		database *Engine
		want     []*api.Repo
	}{
		{
			failure:  false,
			name:     "postgres with name",
			database: _postgres,
			sort:     "name",
			want:     []*api.Repo{_repoOne, _repoTwo},
		},
		{
			failure:  false,
			name:     "postgres with latest",
			database: _postgres,
			sort:     "latest",
			want:     []*api.Repo{_repoOne, _repoTwo},
		},
		{
			failure:  false,
			name:     "sqlite with name",
			database: _sqlite,
			sort:     "name",
			want:     []*api.Repo{_repoOne, _repoTwo},
		},
		{
			failure:  false,
			name:     "sqlite with latest",
			database: _sqlite,
			sort:     "latest",
			want:     []*api.Repo{_repoOne, _repoTwo},
		},
	}

	filters := map[string]interface{}{}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListReposForOrg(context.TODO(), "foo", test.sort, filters, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListReposForOrg for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListReposForOrg for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("ListReposForOrg for %s mismatch (-want +got):\n%s", test.name, diff)
			}
		})
	}
}
