// SPDX-License-Identifier: Apache-2.0

package artifact

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestArtifact_Engine_Count(t *testing.T) {
	// setup types
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
	ctx := context.TODO()

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query for the artifacts table
	_mock.ExpectQuery(`SELECT count(*) FROM "artifacts"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateArtifact(ctx, _artifact)
	if err != nil {
		t.Errorf("unable to create artifact for sqlite: %v", err)
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
			want:     1,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     1,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.CountArtifacts(ctx)
			if test.failure {
				if err == nil {
					t.Errorf("Count for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("Count for %s returned err: %v", test.name, err)
			}

			if got != test.want {
				t.Errorf("Count for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
