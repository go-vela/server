// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPipeline_Engine_DeletePipeline(t *testing.T) {
	// setup types
	_pipeline := testPipeline()
	_pipeline.SetID(1)
	_pipeline.SetRepoID(1)
	_pipeline.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	_pipeline.SetRef("refs/heads/master")
	_pipeline.SetType("yaml")
	_pipeline.SetVersion("1")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`DELETE FROM "pipelines" WHERE "pipelines"."id" = $1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreatePipeline(context.TODO(), _pipeline)
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
			err = test.database.DeletePipeline(context.TODO(), _pipeline)

			if test.failure {
				if err == nil {
					t.Errorf("DeletePipeline for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("DeletePipeline for %s returned err: %v", test.name, err)
			}
		})
	}
}
