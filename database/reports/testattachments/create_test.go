package testattachments

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/server/database/testutils"
)

func TestEngine_CreateTestReportAttachment(t *testing.T) {

	_testReportAttachment := testutils.APITestReportAttachment()
	_testReportAttachment.SetID(1)
	_testReportAttachment.SetTestReportID(1)
	_testReportAttachment.SetFileName("foo")
	_testReportAttachment.SetObjectPath("foo/bar")
	_testReportAttachment.SetFileSize(1)
	_testReportAttachment.SetFileType("xml")
	_testReportAttachment.SetPresignedUrl("foobar")
	_testReportAttachment.SetCreated(1)

	_testReport := testutils.APITestReport()
	_testReport.SetID(1)
	// _testReport.ID.Valid = true

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "testattachments" ("test_report_id","file_name","object_path","file_size","file_type","presigned_url","created","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`).
		WithArgs(1, "foo", "foo/bar", 1, "xml", "foobar", 1, 1).
		WillReturnRows(_rows)

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
			got, err := test.database.CreateTestReportAttachment(context.TODO(), _testReportAttachment)

			if test.failure {
				if err == nil {
					t.Errorf("Create for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("Create for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _testReportAttachment) {
				t.Errorf("Create for %s returned %v, want %v", test.name, got, _testReportAttachment)
			}
		})
	}
}
