// SPDX-License-Identifier: Apache-2.0

package testattachments

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestTestReportsAttachments_Engine_Count(t *testing.T) {
	// setup types
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
	_mock.ExpectQuery(`SELECT count(*) FROM "testattachments"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateTestReportAttachment(ctx, _testReportAttachment)
	if err != nil {
		t.Errorf("unable to create test report attachment for sqlite: %v", err)
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
			got, err := test.database.Count(ctx)
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
