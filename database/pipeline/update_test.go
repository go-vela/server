// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestPipeline_Engine_UpdatePipeline(t *testing.T) {
	// setup types
	_repo := testutils.APIRepo()
	_repo.SetID(1)

	_pipeline := testutils.APIPipeline()
	_pipeline.SetID(1)
	_pipeline.SetRepo(_repo)
	_pipeline.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	_pipeline.SetRef("refs/heads/main")
	_pipeline.SetType("yaml")
	_pipeline.SetVersion("1")
	_pipeline.SetData([]byte{})

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "pipelines"
SET "repo_id"=$1,"commit"=$2,"flavor"=$3,"platform"=$4,"ref"=$5,"type"=$6,"version"=$7,"external_secrets"=$8,"internal_secrets"=$9,"services"=$10,"stages"=$11,"steps"=$12,"templates"=$13,"warnings"=$14,"data"=$15
WHERE "id" = $16`).
		WithArgs(1, "48afb5bdc41ad69bf22588491333f7cf71135163", nil, nil, "refs/heads/main", "yaml", "1", false, false, false, false, false, false, nil, AnyArgument{}, 1).
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
		database *Engine
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
			got, err := test.database.UpdatePipeline(context.TODO(), _pipeline)

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
