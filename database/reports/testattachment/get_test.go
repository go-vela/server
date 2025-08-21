// SPDX-License-Identifier: Apache-2.0

package testattachment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestEngine_GetTestAttachment(t *testing.T) {
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
	_mock.ExpectQuery(`SELECT * FROM "testattachments" WHERE id = $1 LIMIT $2`).WithArgs(1, 1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateTestAttachment(context.TODO(), _testAttachment)
	if err != nil {
		t.Errorf("unable to create test attachment for sqlite: %v", err)
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
			got, err := test.database.GetTestAttachment(context.TODO(), 1)

			if test.failure {
				if err == nil {
					t.Errorf("GetTestAttachment for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetTestAttachment for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetTestAttachment mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
