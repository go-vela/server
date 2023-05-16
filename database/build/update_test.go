// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestBuild_Engine_UpdateBuild(t *testing.T) {
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
	_build.SetPreviousName("oldName")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "builds"
SET "user_id"=$1,"hash"=$2,"org"=$3,"name"=$4,"full_name"=$5,"link"=$6,"clone"=$7,"branch"=$8,"topics"=$9,"build_limit"=$10,"timeout"=$11,"counter"=$12,"visibility"=$13,"private"=$14,"trusted"=$15,"active"=$16,"allow_pull"=$17,"allow_push"=$18,"allow_deploy"=$19,"allow_tag"=$20,"allow_comment"=$21,"pipeline_type"=$22,"previous_name"=$23
WHERE "id" = $24`).
		WithArgs(1, AnyArgument{}, "foo", "bar", "foo/bar", nil, nil, nil, AnyArgument{}, AnyArgument{}, AnyArgument{}, AnyArgument{}, "public", false, false, false, false, false, false, false, false, "yaml", "oldName", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

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
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = test.database.UpdateBuild(_build)

			if test.failure {
				if err == nil {
					t.Errorf("UpdateBuild for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateBuild for %s returned err: %v", test.name, err)
			}
		})
	}
}
