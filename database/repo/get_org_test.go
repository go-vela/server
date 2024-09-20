// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
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
	_rows := sqlmock.NewRows(
		[]string{"id", "user_id", "hash", "org", "name", "full_name", "link", "clone", "branch", "topics", "build_limit", "timeout", "counter", "visibility", "private", "trusted", "active", "allow_events", "pipeline_type", "previous_name", "approve_build"}).
		AddRow(1, 1, "baz", "foo", "bar", "foo/bar", "", "", "", "{}", 0, 0, 0, "public", false, false, false, 1, "yaml", "", "")

	_userRows := sqlmock.NewRows(
		[]string{"id", "name", "token", "hash", "active", "admin"}).
		AddRow(1, "foo", "bar", "baz", false, false)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE full_name = $1 LIMIT $2`).WithArgs("foo/bar", 1).WillReturnRows(_rows)
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
		database *engine
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
			got, err := test.database.GetRepoForOrg(context.TODO(), "foo/bar")

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
