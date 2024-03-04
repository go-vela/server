// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestPipeline_Engine_GetPipelineForRepo(t *testing.T) {
	// setup types
	_pipeline := testPipeline()
	_pipeline.SetID(1)
	_pipeline.SetRepoID(1)
	_pipeline.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	_pipeline.SetRef("refs/heads/main")
	_pipeline.SetType("yaml")
	_pipeline.SetVersion("1")
	_pipeline.SetData([]byte("foo"))

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "commit", "flavor", "platform", "ref", "type", "version", "services", "stages", "steps", "templates", "data"}).
		AddRow(1, 1, "48afb5bdc41ad69bf22588491333f7cf71135163", "", "", "refs/heads/main", "yaml", "1", false, false, false, false, []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "pipelines" WHERE repo_id = $1 AND "commit" = $2 LIMIT $3`).WithArgs(1, "48afb5bdc41ad69bf22588491333f7cf71135163", 1).WillReturnRows(_rows)

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
		want     *library.Pipeline
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _pipeline,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _pipeline,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetPipelineForRepo(context.TODO(), "48afb5bdc41ad69bf22588491333f7cf71135163", &library.Repo{ID: _pipeline.RepoID})

			if test.failure {
				if err == nil {
					t.Errorf("GetPipelineForRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetPipelineForRepo for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetPipelineForRepo for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
