// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestRepo_Engine_GetRepoForOrg(t *testing.T) {
	// setup types
	_repo := testutils.APIRepo()
	_repo.SetID(1)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")
	_repo.SetPipelineType("yaml")
	_repo.SetTopics([]string{})
	_repo.SetAllowEvents(api.NewEventsFromMask(1))

	_owner := testutils.APIUser().Crop()
	_owner.SetID(1)
	_owner.SetName("foo")
	_owner.SetToken("bar")

	_repo.SetOwner(_owner)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.RepoFromAPI(_repo)})

	_userRows := testutils.CreateMockRows([]any{*types.UserFromAPI(_owner)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE org = $1 AND name = $2 LIMIT $3`).WithArgs("foo", "bar", 1).WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."id" = $1`).WithArgs(1).WillReturnRows(_userRows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateRepo(context.TODO(), _repo)
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
		failure  bool
		name     string
		database *Engine
		want     *api.Repo
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _repo,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _repo,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetRepoForOrg(context.TODO(), "foo", "bar")

			if test.failure {
				if err == nil {
					t.Errorf("GetRepoForOrg for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetRepoForOrg for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(_repo, got); diff != "" {
				t.Errorf("GetRepoForOrg mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
