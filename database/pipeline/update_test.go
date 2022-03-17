// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPipeline_Engine_UpdatePipeline(t *testing.T) {
	// setup types
	_pipeline := testPipeline()
	_pipeline.SetID(1)
	_pipeline.SetRepoID(1)
	_pipeline.SetNumber(1)
	_pipeline.SetRef("refs/heads/master")
	_pipeline.SetType("yaml")
	_pipeline.SetVersion("1")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "pipelines"
SET "repo_id"=$1,"number"=$2,"commit"=$3,"flavor"=$4,"platform"=$5,"ref"=$6,"type"=$7,"version"=$8,"external_secrets"=$9,"internal_secrets"=$10,"services"=$11,"stages"=$12,"steps"=$13,"templates"=$14,"data"=$15
WHERE "id" = $16`).
		WithArgs(1, 1, nil, nil, nil, "refs/heads/master", "yaml", "1", false, false, false, false, false, false, AnyArgument{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreatePipeline(_pipeline)
	if err != nil {
		t.Errorf("unable to create test pipeline for sqlite: %v", err)
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
			err = test.database.UpdatePipeline(_pipeline)

			if test.failure {
				if err == nil {
					t.Errorf("UpdatePipeline for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdatePipeline for %s returned err: %v", test.name, err)
			}
		})
	}
}
