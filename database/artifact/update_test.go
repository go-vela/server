// SPDX-License-Identifier: Apache-2.0

package artifact

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
)

func TestArtifact_Engine_Update(t *testing.T) {
	// setup types
	_owner := testutils.APIUser()
	_owner.SetID(1)
	_owner.SetName("foo")
	_owner.SetToken("bar")

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
	_mock.ExpectExec(`UPDATE "artifacts" SET "build_id"=$1,"file_name"=$2,"object_path"=$3,"file_size"=$4,"file_type"=$5,"presigned_url"=$6,"created_at"=$7 WHERE "id" = $8`).
		WithArgs(1, "foo", "foo/bar", 1, "xml", "foobar", 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateArtifact(ctx, _artifact)
	if err != nil {
		t.Errorf("unable to update artifact for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     *api.Artifact
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _artifact,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _artifact,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.UpdateArtifact(ctx, _artifact)

			if test.failure {
				if err == nil {
					t.Errorf("Update for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("Update for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetArtifact mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
