// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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
	_pipeline.SetRef("48afb5bdc41ad69bf22588491333f7cf71135163")
	_pipeline.SetType("yaml")
	_pipeline.SetVersion("1")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "pipelines"
SET "repo_id"=$1,"number"=$2,"flavor"=$3,"platform"=$4,"ref"=$5,"type"=$6,"version"=$7,"services"=$8,"stages"=$9,"steps"=$10,"templates"=$11,"data"=$12
WHERE "id" = $13`).
		WithArgs(1, 1, nil, nil, "48afb5bdc41ad69bf22588491333f7cf71135163", "yaml", "1", false, false, false, false, AnyArgument{}, 1).
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
			name:     "sqlite",
			database: _sqlite,
		},
	}

	// run tests
	for _, test := range tests {
		err = test.database.UpdatePipeline(_pipeline)

		if test.failure {
			if err == nil {
				t.Errorf("UpdatePipeline for %s should have returned err", test.name)
			}

			continue
		}

		if err != nil {
			t.Errorf("UpdatePipeline for %s returned err: %v", test.name, err)
		}
	}
}
