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

func TestTestAttachment_Engine_ListTestAttachments(t *testing.T) {
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

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.TestAttachmentFromAPI(_testAttachment)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "testattachments" ORDER BY created_at DESC`).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	// Create necessary SQLite tables for relationship testing
	err := _sqlite.client.AutoMigrate(&types.TestReport{})
	if err != nil {
		t.Errorf("unable to create tables for sqlite: %v", err)
	}

	// Create the test report in sqlite
	_, err = _sqlite.CreateTestAttachment(ctx, _testAttachment)
	if err != nil {
		t.Errorf("unable to create test report for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     []*api.TestAttachment
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*api.TestAttachment{_testAttachment},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.TestAttachment{_testAttachment},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListTestAttachments(ctx)

			if test.failure {
				if err == nil {
					t.Errorf("ListTestAttachments should have returned err")
				}

				return
			}

			if err != nil {
				t.Errorf("ListTestAttachments returned err: %v", err)
			}

			if len(got) != len(test.want) {
				t.Errorf("ListTestAttachments for %s returned %d test attachments, want %d", test.name, len(got), len(test.want))
				return
			}

			if len(got) > 0 {
				// Check report fields
				if !reflect.DeepEqual(got[0].GetID(), test.want[0].GetID()) ||
					!reflect.DeepEqual(got[0].GetCreatedAt(), test.want[0].GetCreatedAt()) {
					t.Errorf("ListTestAttachments for %s returned unexpected test attachments values: got %v, want %v",
						test.name, got[0], test.want[0])
				}
			}
		})
	}
}
