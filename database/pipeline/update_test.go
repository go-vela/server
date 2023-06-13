// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPipeline_Engine_UpdatePipeline(t *testing.T) {
	// setup types
	_pipeline := testPipeline()
	_pipeline.SetID(1)
	_pipeline.SetRepoID(1)
	_pipeline.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	_pipeline.SetRef("refs/heads/master")
	_pipeline.SetType("yaml")
	_pipeline.SetVersion("1")
	_pipeline.SetData([]byte{})

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "pipelines"
SET "repo_id"=$1,"commit"=$2,"flavor"=$3,"platform"=$4,"ref"=$5,"type"=$6,"version"=$7,"external_secrets"=$8,"internal_secrets"=$9,"services"=$10,"stages"=$11,"steps"=$12,"templates"=$13,"data"=$14
WHERE "id" = $15`).
		WithArgs(1, "48afb5bdc41ad69bf22588491333f7cf71135163", nil, nil, "refs/heads/master", "yaml", "1", false, false, false, false, false, false, AnyArgument{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreatePipeline(_pipeline)
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
			got, err := test.database.UpdatePipeline(_pipeline)

			if test.failure {
				if err == nil {
					t.Errorf("UpdatePipeline for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdatePipeline for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _pipeline) {
				t.Errorf("UpdatePipeline for %s returned %s, want %s", test.name, got, _pipeline)
			}
		})
	}
}
