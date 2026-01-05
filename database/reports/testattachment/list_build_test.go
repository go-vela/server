// SPDX-License-Identifier: Apache-2.0

package testattachment

import (
	"context"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestTestAttachment_Engine_ListTestAttachmentsByBuildID(t *testing.T) {
	// setup types
	ctx := context.Background()
	_testAttachment := testutils.APITestAttachment()
	_testAttachment.SetID(1)
	_testAttachment.SetTestReportID(1)
	_testAttachment.SetFileName("foo")
	_testAttachment.SetObjectPath("foo/bar")
	_testAttachment.SetFileSize(1)
	_testAttachment.SetFileType("xml")
	_testAttachment.SetPresignedURL("foobar")
	_testAttachment.SetCreatedAt(1)

	_testReport := testutils.APITestReport()
	_testReport.SetID(1)
	_testReport.SetBuildID(1)
	_testReport.SetCreatedAt(1)

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	// create SQLite tables for relationship testing with correct table names
	err := _sqlite.client.Exec(`
		CREATE TABLE IF NOT EXISTS test_reports (
			id INTEGER PRIMARY KEY,
			build_id INTEGER,
			created_at INTEGER
		)`).Error
	if err != nil {
		t.Errorf("unable to create test_reports table for sqlite: %v", err)
	}

	err = _sqlite.client.Exec(`
		CREATE TABLE IF NOT EXISTS testattachments (
			id INTEGER PRIMARY KEY,
			test_report_id INTEGER,
			file_name TEXT,
			object_path TEXT,
			file_size INTEGER,
			file_type TEXT,
			presigned_url TEXT,
			created_at INTEGER
		)`).Error
	if err != nil {
		t.Errorf("unable to create testattachments table for sqlite: %v", err)
	}

	// create the test report in sqlite
	err = _sqlite.client.Create(types.TestReportFromAPI(_testReport)).Error
	if err != nil {
		t.Errorf("unable to create test report for sqlite: %v", err)
	}

	// create the test attachment in sqlite
	_, err = _sqlite.CreateTestAttachment(ctx, _testAttachment)
	if err != nil {
		t.Errorf("unable to create test attachment for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		buildID  int64
		want     []*api.TestAttachment
	}{
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			buildID:  1,
			want:     []*api.TestAttachment{_testAttachment},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListTestAttachmentsByBuildID(ctx, test.buildID)

			if test.failure {
				if err == nil {
					t.Errorf("ListTestAttachmentsByBuildID should have returned err")
				}

				return
			}

			if err != nil {
				t.Errorf("ListTestAttachmentsByBuildID returned err: %v", err)
			}

			if len(got) != len(test.want) {
				t.Errorf("ListTestAttachmentsByBuildID for %s returned %d test attachments, want %d", test.name, len(got), len(test.want))
				return
			}

			if len(got) > 0 {
				// check attachment fields
				if !reflect.DeepEqual(got[0].GetID(), test.want[0].GetID()) ||
					!reflect.DeepEqual(got[0].GetCreatedAt(), test.want[0].GetCreatedAt()) {
					t.Errorf("ListTestAttachmentsByBuildID for %s returned %v, want %v", test.name, got[0], test.want[0])
				}
			}
		})
	}
}
