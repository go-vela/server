// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRepo_Engine_UpdateRepo(t *testing.T) {
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
	_repo.SetPreviousName("oldName")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "repos"
SET "user_id"=$1,"hash"=$2,"org"=$3,"name"=$4,"full_name"=$5,"link"=$6,"clone"=$7,"branch"=$8,"topics"=$9,"build_limit"=$10,"timeout"=$11,"counter"=$12,"visibility"=$13,"private"=$14,"trusted"=$15,"active"=$16,"allow_pull"=$17,"allow_push"=$18,"allow_deploy"=$19,"allow_tag"=$20,"allow_comment"=$21,"pipeline_type"=$22,"previous_name"=$23
WHERE "id" = $24`).
		WithArgs(1, AnyArgument{}, "foo", "bar", "foo/bar", nil, nil, nil, AnyArgument{}, AnyArgument{}, AnyArgument{}, AnyArgument{}, "public", false, false, false, false, false, false, false, false, "yaml", "oldName", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

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
			err = test.database.UpdateRepo(_repo)

			if test.failure {
				if err == nil {
					t.Errorf("UpdateRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateRepo for %s returned err: %v", test.name, err)
			}
		})
	}
}
