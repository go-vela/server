// SPDX-License-Identifier: Apache-2.0

package artifact

import (
	"context"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestArtifact_Engine_ListArtifacts(t *testing.T) {
	// setup types
	ctx := context.Background()
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
	_mock.ExpectQuery(`SELECT * FROM "artifacts" ORDER BY created_at DESC`).WillReturnRows(_rows)

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	// Create the artifact in sqlite
	_, err := _sqlite.CreateArtifact(ctx, _artifact)
	if err != nil {
		t.Errorf("unable to create artifact for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     []*api.Artifact
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*api.Artifact{_artifact},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.Artifact{_artifact},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListArtifacts(ctx)

			if test.failure {
				if err == nil {
					t.Errorf("ListArtifacts should have returned err")
				}

				return
			}

			if err != nil {
				t.Errorf("ListArtifacts returned err: %v", err)
			}

			if len(got) != len(test.want) {
				t.Errorf("ListArtifacts for %s returned %d artifacts, want %d", test.name, len(got), len(test.want))
				return
			}

			if len(got) > 0 {
				// Check report fields
				if !reflect.DeepEqual(got[0].GetID(), test.want[0].GetID()) ||
					!reflect.DeepEqual(got[0].GetCreatedAt(), test.want[0].GetCreatedAt()) {
					t.Errorf("ListArtifacts for %s returned unexpected artifacts values: got %v, want %v",
						test.name, got[0], test.want[0])
				}
			}
		})
	}
}
