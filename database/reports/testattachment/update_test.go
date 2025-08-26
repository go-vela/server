// SPDX-License-Identifier: Apache-2.0

package testattachment

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
)

func TestTestAttachment_Engine_Update(t *testing.T) {
	// setup types
	_owner := testutils.APIUser()
	_owner.SetID(1)
	_owner.SetName("foo")
	_owner.SetToken("bar")

	_testAttachment := testutils.APITestAttachment()
	_testAttachment.SetID(1)
	_testAttachment.SetTestReportID(1)
	_testAttachment.SetFileName("foo")
	_testAttachment.SetObjectPath("foo/bar")
	_testAttachment.SetFileSize(1)
	_testAttachment.SetFileType("xml")
	_testAttachment.SetPresignedURL("foobar")
	_testAttachment.SetCreatedAt(1)

	_postgres, _mock := testPostgres(t)
	ctx := context.TODO()

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query for the test_reports table
	_mock.ExpectExec(`UPDATE "testattachments" SET "test_report_id"=$1,"file_name"=$2,"object_path"=$3,"file_size"=$4,"file_type"=$5,"presigned_url"=$6,"created_at"=$7 WHERE "id" = $8`).
		WithArgs(1, "foo", "foo/bar", 1, "xml", "foobar", 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateTestAttachment(ctx, _testAttachment)
	if err != nil {
		t.Errorf("unable to update test attachment for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     *api.TestAttachment
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _testAttachment,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _testAttachment,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.UpdateTestAttachment(ctx, _testAttachment)

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
				t.Errorf("GetTestAttachment mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
