// SPDX-License-Identifier: Apache-2.0

package testreports

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestTestReports_Engine_Count(t *testing.T) {
	// setup types
	_testReport := testutils.APITestReport()
	_testReport.SetID(1)
	_testReport.SetBuildID(1)
	_testReport.SetCreated(1)

	_postgres, _mock := testPostgres(t)
	ctx := context.TODO()
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query for the test_reports table
	_mock.ExpectQuery(`SELECT count(*) FROM "testreports"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.Create(ctx, _testReport)
	if err != nil {
		t.Errorf("unable to create test report for sqlite: %v", err)
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
