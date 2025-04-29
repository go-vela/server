// SPDX-License-Identifier: Apache-2.0

package testattachments

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
)

func TestTestReportsAttachment_Engine_Update(t *testing.T) {
	// setup types
	_owner := testutils.APIUser()
	_owner.SetID(1)
	_owner.SetName("foo")
	_owner.SetToken("bar")

	_testReportAttachment := testutils.APITestReportAttachment()
	_testReportAttachment.SetID(1)
	_testReportAttachment.SetTestReportID(1)
	_testReportAttachment.SetFileName("foo")
	_testReportAttachment.SetObjectPath("foo/bar")
	_testReportAttachment.SetFileSize(1)
	_testReportAttachment.SetFileType("xml")
	_testReportAttachment.SetPresignedUrl("foobar")
	_testReportAttachment.SetCreated(1)

	_postgres, _mock := testPostgres(t)
	ctx := context.TODO()
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query for the test_reports table
	_mock.ExpectExec(`UPDATE "testattachments" SET "test_report_id"=$1,"file_name"=$2,"object_path"=$3,"file_size"=$4,"file_type"=$5,"presigned_url"=$6,"created"=$7 WHERE "id" = $8`).
		WithArgs(1, "foo", "foo/bar", 1, "xml", "foobar", 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateTestReportAttachment(ctx, _testReportAttachment)
	if err != nil {
		t.Errorf("unable to update test report attachment for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     *api.TestReportAttachments
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _testReportAttachment,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _testReportAttachment,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.Update(ctx, _testReportAttachment)

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
				t.Errorf("GetTestReportAttachment mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
