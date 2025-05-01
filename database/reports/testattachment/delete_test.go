// SPDX-License-Identifier: Apache-2.0

package testattachment

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestTestAttachment_Engine_Delete(t *testing.T) {
	// setup types
	_testAttachment := testutils.APITestAttachment()
	_testAttachment.SetID(1)
	_testAttachment.SetTestReportID(1)
	_testAttachment.SetFileName("foo")
	_testAttachment.SetObjectPath("foo/bar")
	_testAttachment.SetFileSize(1)
	_testAttachment.SetFileType("xml")
	_testAttachment.SetPresignedUrl("foobar")
	_testAttachment.SetCreatedAt(1)

	_postgres, _mock := testPostgres(t)
	ctx := context.TODO()
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	_mock.ExpectExec(`DELETE FROM "testattachments" WHERE "testattachments"."id" = $1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

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
			err := test.database.DeleteTestAttachment(ctx, _testAttachment)

			if test.failure {
				if err == nil {
					t.Errorf("Delete for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("Delete for %s returned err: %v", test.name, err)
			}
		})
	}
}
