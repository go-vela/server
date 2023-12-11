// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestRepo_Engine_GetRepoForOrg(t *testing.T) {
	// setup types
	_repo := testRepo()
	_repo.SetID(1)
	_repo.SetUserID(1)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")
	_repo.SetPipelineType("yaml")
	_repo.SetTopics([]string{})
	_repo.SetAllowEvents(library.NewEventsFromMask(1))

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "user_id", "hash", "org", "name", "full_name", "link", "clone", "branch", "topics", "build_limit", "timeout", "counter", "visibility", "private", "trusted", "active", "allow_pull", "allow_push", "allow_deploy", "allow_tag", "allow_comment", "allow_events", "pipeline_type", "previous_name", "approve_build"}).
		AddRow(1, 1, "baz", "foo", "bar", "foo/bar", "", "", "", "{}", 0, 0, 0, "public", false, false, false, false, false, false, false, false, 1, "yaml", "", "")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE org = $1 AND name = $2 LIMIT 1`).WithArgs("foo", "bar").WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateRepo(context.TODO(), _repo)
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     *library.Repo
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

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetRepoForOrg for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
