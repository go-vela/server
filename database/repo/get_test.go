// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestRepo_Engine_GetRepo(t *testing.T) {
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

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "user_id", "hash", "org", "name", "full_name", "link", "clone", "branch", "build_limit", "timeout", "counter", "visibility", "private", "trusted", "active", "allow_pull", "allow_pull_fork", "allow_push", "allow_deploy", "allow_tag", "allow_comment", "pipeline_type", "previous_name"}).
		AddRow(1, 1, "baz", "foo", "bar", "foo/bar", "", "", "", 0, 0, 0, "public", false, false, false, false, false, false, false, false, false, "yaml", "")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE id = $1 LIMIT 1`).WithArgs(1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateRepo(_repo)
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
			got, err := test.database.GetRepo(1)

			if test.failure {
				if err == nil {
					t.Errorf("GetRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetRepo for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetRepo for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
