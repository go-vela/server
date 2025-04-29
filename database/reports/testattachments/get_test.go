// SPDX-License-Identifier: Apache-2.0

package testattachments

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestEngine_GetTestReportAttachment(t *testing.T) {
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
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.TestReportAttachmentFromAPI(_testReportAttachment)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "testattachments" WHERE id = $1 LIMIT $2`).WithArgs(1, 1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateTestReportAttachment(context.TODO(), _testReportAttachment)
	if err != nil {
		t.Errorf("unable to create test report attachment for sqlite: %v", err)
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
			got, err := test.database.GetTestReportAttachment(context.TODO(), 1)

			if test.failure {
				if err == nil {
					t.Errorf("GetTestReportAttachment for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetTestReportAttachment for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetTestReportAttachment mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
