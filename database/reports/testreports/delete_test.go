// SPDX-License-Identifier: Apache-2.0

package testreports

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestTestReports_Engine_Delete(t *testing.T) {
	// setup types
	_report := testutils.APITestReport()
	_report.SetID(1)
	_report.SetBuildID(1)
	_report.SetCreated(1)

	_postgres, _mock := testPostgres(t)
	ctx := context.TODO()
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	_mock.ExpectExec(`DELETE FROM "testreports" WHERE "testreports"."id" = $1`).
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
			err := test.database.DeleteByID(ctx, _report)

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
