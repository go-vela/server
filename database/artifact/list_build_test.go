// SPDX-License-Identifier: Apache-2.0

package artifact

import (
	"context"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
)

func TestArtifact_Engine_ListArtifactsByBuildID(t *testing.T) {
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

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	// create SQLite tables for relationship testing with correct table names
	err := _sqlite.client.Exec(`
		CREATE TABLE IF NOT EXISTS artifacts (
			id INTEGER PRIMARY KEY,
			build_id INTEGER,
			created_at INTEGER
		)`).Error
	if err != nil {
		t.Errorf("unable to create artifacts table for sqlite: %v", err)
	}

	err = _sqlite.client.Exec(`
		CREATE TABLE IF NOT EXISTS artifacts (
			id INTEGER PRIMARY KEY,
			build_id INTEGER,
			file_name TEXT,
			object_path TEXT,
			file_size INTEGER,
			file_type TEXT,
			presigned_url TEXT,
			created_at INTEGER
		)`).Error
	if err != nil {
		t.Errorf("unable to create artifacts table for sqlite: %v", err)
	}

	// create the artifact in sqlite
	_, err = _sqlite.CreateArtifact(ctx, _artifact)
	if err != nil {
		t.Errorf("unable to create artifact for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		buildID  int64
		want     []*api.Artifact
	}{
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			buildID:  1,
			want:     []*api.Artifact{_artifact},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListArtifactsByBuildID(ctx, test.buildID)

			if test.failure {
				if err == nil {
					t.Errorf("ListArtifactsByBuildID should have returned err")
				}

				return
			}

			if err != nil {
				t.Errorf("ListArtifactsByBuildID returned err: %v", err)
			}

			if len(got) != len(test.want) {
				t.Errorf("ListArtifactsByBuildID for %s returned %d artifacts, want %d", test.name, len(got), len(test.want))
				return
			}

			if len(got) > 0 {
				// check artifact fields
				if !reflect.DeepEqual(got[0].GetID(), test.want[0].GetID()) ||
					!reflect.DeepEqual(got[0].GetCreatedAt(), test.want[0].GetCreatedAt()) {
					t.Errorf("ListArtifactsByBuildID for %s returned %v, want %v", test.name, got[0], test.want[0])
				}
			}
		})
	}
}
