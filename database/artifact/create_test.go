// SPDX-License-Identifier: Apache-2.0

package artifact

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestEngine_CreateArtifact(t *testing.T) {
	_artifact := testutils.APIArtifact()
	_artifact.SetID(1)
	_artifact.SetBuildID(1)
	_artifact.SetFileName("foo")
	_artifact.SetObjectPath("foo/bar")
	_artifact.SetFileSize(1)
	_artifact.SetFileType("xml")
	_artifact.SetPresignedURL("foobar")
	_artifact.SetCreatedAt(1)

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "artifacts" ("build_id","file_name","object_path","file_size","file_type","presigned_url","created_at","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`).
		WithArgs(1, "foo", "foo/bar", 1, "xml", "foobar", 1, 1).
		WillReturnRows(_rows)

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

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
			got, err := test.database.CreateArtifact(context.TODO(), _artifact)

			if test.failure {
				if err == nil {
					t.Errorf("Create for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("Create for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _artifact) {
				t.Errorf("Create for %s returned %v, want %v", test.name, got, _artifact)
			}
		})
	}
}
