// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestPipeline_Engine_CountPipelinesForRepo(t *testing.T) {
	// setup types
	_repo := testutils.APIRepo()
	_repo.SetID(1)

	_pipelineOne := testutils.APIPipeline()
	_pipelineOne.SetID(1)
	_pipelineOne.SetRepo(_repo)
	_pipelineOne.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	_pipelineOne.SetRef("refs/heads/main")
	_pipelineOne.SetType("yaml")
	_pipelineOne.SetVersion("1")

	_pipelineTwo := testutils.APIPipeline()
	_pipelineTwo.SetID(2)
	_pipelineTwo.SetRepo(_repo)
	_pipelineTwo.SetCommit("a49aaf4afae6431a79239c95247a2b169fd9f067")
	_pipelineTwo.SetRef("refs/heads/main")
	_pipelineTwo.SetType("yaml")
	_pipelineTwo.SetVersion("1")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "pipelines" WHERE repo_id = $1`).WithArgs(1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreatePipeline(context.TODO(), _pipelineOne)
	if err != nil {
		t.Errorf("unable to create test pipeline for sqlite: %v", err)
	}

	_, err = _sqlite.CreatePipeline(context.TODO(), _pipelineTwo)
	if err != nil {
		t.Errorf("unable to create test pipeline for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     int64
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     2,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     2,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.CountPipelinesForRepo(context.TODO(), _repo)

			if test.failure {
				if err == nil {
					t.Errorf("CountPipelinesForRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CountPipelinesForRepo for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("CountPipelinesForRepo for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
