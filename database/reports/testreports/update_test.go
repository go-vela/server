// SPDX-License-Identifier: Apache-2.0

package testreports

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestTestReports_Engine_Update(t *testing.T) {
	// setup types
	_owner := testutils.APIUser()
	_owner.SetID(1)
	_owner.SetName("foo")
	_owner.SetToken("bar")

	_testReport := testutils.APITestReport()
	_testReport.SetID(1)
	_testReport.SetBuildID(1)
	_testReport.SetCreated(1)

	_testReportAttachments := testutils.APITestReportAttachment()
	_testReportAttachments.SetID(1)
	_testReportAttachments.SetTestReportID(1)
	_testReportAttachments.SetFilename("test.xml")
	_testReportAttachments.SetFilePath("/path/to/test.xml")
	_testReportAttachments.SetFileSize(1024)
	_testReportAttachments.SetFileType("application/xml")
	_testReportAttachments.SetPresignedUrl("https://example.com/test.xml")
	_testReportAttachments.SetCreated(1)

	_testReport.SetReportAttachments(_testReportAttachments)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query for the test_reports table
	_mock.ExpectBegin()
	_mock.ExpectExec(`UPDATE "test_reports" SET .* WHERE "id" = .`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectCommit()
	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateTestReport(context.TODO(), _testReport)
	if err != nil {
		t.Errorf("unable to create test report for sqlite: %v", err)
	}

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
			got, err := test.database.Update(context.TODO(), _testReport)

			if test.failure {
				if err == nil {
					t.Errorf("Update for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("Update for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got.GetID(), _testReport.GetID()) ||
				!reflect.DeepEqual(got.GetBuildID(), _testReport.GetBuildID()) ||
				!reflect.DeepEqual(got.GetCreated(), _testReport.GetCreated()) {
				t.Errorf("Update for %s returned unexpected report values: got %v, want %v", test.name, got, _testReport)
			}

			// Check attachment values
			//if !reflect.DeepEqual(got.ReportAttachments.GetID(), _testReport.ReportAttachments.GetID()) ||
			//	!reflect.DeepEqual(got.ReportAttachments.GetTestReportID(), _testReport.ReportAttachments.GetTestReportID()) ||
			//	!reflect.DeepEqual(got.ReportAttachments.GetFilename(), _testReport.ReportAttachments.GetFilename()) ||
			//	!reflect.DeepEqual(got.ReportAttachments.GetFilePath(), _testReport.ReportAttachments.GetFilePath()) ||
			//	!reflect.DeepEqual(got.ReportAttachments.GetFileSize(), _testReport.ReportAttachments.GetFileSize()) ||
			//	!reflect.DeepEqual(got.ReportAttachments.GetFileType(), _testReport.ReportAttachments.GetFileType()) ||
			//	!reflect.DeepEqual(got.ReportAttachments.GetPresignedUrl(), _testReport.ReportAttachments.GetPresignedUrl()) ||
			//	!reflect.DeepEqual(got.ReportAttachments.GetCreated(), _testReport.ReportAttachments.GetCreated()) {
			//	t.Errorf("Update for %s returned unexpected attachment values: got %v, want %v",
			//		test.name, got.ReportAttachments, _testReport.ReportAttachments)
			//}
		})
	}
}
