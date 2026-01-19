// SPDX-License-Identifier: Apache-2.0

package artifact

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestEngine_GetArtifact(t *testing.T) {
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
	_rows := testutils.CreateMockRows([]any{*types.ArtifactFromAPI(_artifact)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "artifacts" WHERE id = $1 LIMIT $2`).WithArgs(1, 1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateArtifact(context.TODO(), _artifact)
	if err != nil {
		t.Errorf("unable to create artifact for sqlite: %v", err)
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
			got, err := test.database.GetArtifact(context.TODO(), 1)

			if test.failure {
				if err == nil {
					t.Errorf("GetArtifact for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetArtifact for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetArtifact mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
