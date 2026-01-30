// SPDX-License-Identifier: Apache-2.0

package artifact

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestArtifact_Engine_Delete(t *testing.T) {
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

	_mock.ExpectExec(`DELETE FROM "artifacts" WHERE "artifacts"."id" = $1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

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
			err := test.database.DeleteArtifact(ctx, _artifact)

			if test.failure {
				if err == nil {
					t.Errorf("Delete for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("Delete for %s returned err: %v", test.name, err)
			}
		})
	}
}
