// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPipeline_Engine_CreatePipeline(t *testing.T) {
	// setup types
	_pipeline := testPipeline()
	_pipeline.SetID(1)
	_pipeline.SetRepoID(1)
	_pipeline.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	_pipeline.SetRef("refs/heads/main")
	_pipeline.SetType("yaml")
	_pipeline.SetVersion("1")
	_pipeline.SetData([]byte{})

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "pipelines"
("repo_id","commit","flavor","platform","ref","type","version","external_secrets","internal_secrets","services","stages","steps","templates","data","id")
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) RETURNING "id"`).
		WithArgs(1, "48afb5bdc41ad69bf22588491333f7cf71135163", nil, nil, "refs/heads/main", "yaml", "1", false, false, false, false, false, false, AnyArgument{}, 1).
		WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

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
			got, err := test.database.CreatePipeline(context.TODO(), _pipeline)

			if test.failure {
				if err == nil {
					t.Errorf("CreatePipeline for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreatePipeline for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _pipeline) {
				t.Errorf("CreatePipeline for %s returned %s, want %s", test.name, got, _pipeline)
			}
		})
	}
}
