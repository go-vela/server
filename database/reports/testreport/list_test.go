// SPDX-License-Identifier: Apache-2.0

package testreport

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestTestReport_Engine_ListTestReports(t *testing.T) {
	// setup types
	ctx := context.Background()
	_testReport := testutils.APITestReport()
	_testReport.SetID(1)
	_testReport.SetBuildID(1)
	_testReport.SetCreatedAt(1)

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query for the test_reports table
	_rows := sqlmock.NewRows([]string{"id", "build_id", "created_at"}).
		AddRow(1, 1, 1)
	_mock.ExpectQuery(`SELECT * FROM "testreports" ORDER BY created_at DESC`).
		WillReturnRows(_rows)

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	// Create necessary SQLite tables for relationship testing
	err := _sqlite.client.AutoMigrate(&types.TestReport{})
	if err != nil {
		t.Errorf("unable to create tables for sqlite: %v", err)
	}

	// Create the test report in sqlite
	_, err = _sqlite.CreateTestReport(ctx, _testReport)
	if err != nil {
		t.Errorf("unable to create test report for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     []*api.TestReport
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*api.TestReport{_testReport},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.TestReport{_testReport},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListTestReports(ctx)

			if test.failure {
				if err == nil {
					t.Errorf("ListTestReports should have returned err")
				}

				return
			}

			if err != nil {
				t.Errorf("ListTestReports returned err: %v", err)
			}

			if len(got) != len(test.want) {
				t.Errorf("ListTestReports for %s returned %d reports, want %d", test.name, len(got), len(test.want))
				return
			}

			if len(got) > 0 {
				// Check report fields
				if !reflect.DeepEqual(got[0].GetID(), test.want[0].GetID()) ||
					!reflect.DeepEqual(got[0].GetBuildID(), test.want[0].GetBuildID()) ||
					!reflect.DeepEqual(got[0].GetCreatedAt(), test.want[0].GetCreatedAt()) {
					t.Errorf("ListTestReports for %s returned unexpected report values: got %v, want %v",
						test.name, got[0], test.want[0])
				}
			}
		})
	}
}
