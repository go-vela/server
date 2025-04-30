// SPDX-License-Identifier: Apache-2.0

package testattachments

import (
	"context"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestTestAttachments_Engine_ListTestAttachments(t *testing.T) {
	// setup types
	ctx := context.Background()
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
	_mock.ExpectQuery(`SELECT * FROM "testattachments" ORDER BY created DESC`).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	// Create necessary SQLite tables for relationship testing
	err := _sqlite.client.AutoMigrate(&types.TestReport{})
	if err != nil {
		t.Errorf("unable to create tables for sqlite: %v", err)
	}

	// Create the test report in sqlite
	_, err = _sqlite.CreateTestReportAttachment(ctx, _testReportAttachment)
	if err != nil {
		t.Errorf("unable to create test report for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     []*api.TestReportAttachments
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*api.TestReportAttachments{_testReportAttachment},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.TestReportAttachments{_testReportAttachment},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.List(ctx)

			if test.failure {
				if err == nil {
					t.Errorf("List should have returned err")
				}

				return
			}

			if err != nil {
				t.Errorf("List returned err: %v", err)
			}

			if len(got) != len(test.want) {
				t.Errorf("List for %s returned %d reports, want %d", test.name, len(got), len(test.want))
				return
			}

			if len(got) > 0 {
				// Check report fields
				if !reflect.DeepEqual(got[0].GetID(), test.want[0].GetID()) ||
					!reflect.DeepEqual(got[0].GetCreated(), test.want[0].GetCreated()) {
					t.Errorf("List for %s returned unexpected report values: got %v, want %v",
						test.name, got[0], test.want[0])
				}
			}
		})
	}
}
