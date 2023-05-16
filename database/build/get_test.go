// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestBuild_Engine_GetBuild(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetUserID(1)
	_build.SetHash("baz")
	_build.SetOrg("foo")
	_build.SetName("bar")
	_build.SetFullName("foo/bar")
	_build.SetVisibility("public")
	_build.SetPipelineType("yaml")
	_build.SetTopics([]string{})

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "user_id", "hash", "org", "name", "full_name", "link", "clone", "branch", "topics", "build_limit", "timeout", "counter", "visibility", "private", "trusted", "active", "allow_pull", "allow_push", "allow_deploy", "allow_tag", "allow_comment", "pipeline_type", "previous_name"}).
		AddRow(1, 1, "baz", "foo", "bar", "foo/bar", "", "", "", "{}", 0, 0, 0, "public", false, false, false, false, false, false, false, false, "yaml", "")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "builds" WHERE id = $1 LIMIT 1`).WithArgs(1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateBuild(_build)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     *library.Build
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _build,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _build,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetBuild(1)

			if test.failure {
				if err == nil {
					t.Errorf("GetBuild for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetBuild for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetBuild for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
