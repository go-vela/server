// SPDX-License-Identifier: Apache-2.0

package testreport

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestTestReport_Engine_CountByBuild(t *testing.T) {
	// setup types
	_testreport := testutils.APITestReport()
	_testreport.SetID(1)
	_testreport.SetBuildID(1)
	_testreport.SetCreatedAt(1)

	_build := testutils.APIBuild()
	_build.SetID(1)

	_postgres, _mock := testPostgres(t)
	ctx := context.TODO()
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query for the test_reports table
	_mock.ExpectQuery(`SELECT count(*) FROM "testreports" WHERE build_id = $1 AND created_at < $2 AND created_at > $3`).
		WithArgs(1, AnyArgument{}, 0).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateTestReport(ctx, _testreport)
	if err != nil {
		t.Errorf("unable to create test report: %v", err)
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

	filters := map[string]interface{}{}
	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.CountTestReportsByBuild(ctx, _build, filters, time.Now().Unix(), 0)
			if test.failure {
				if err == nil {
					t.Errorf("CountTestReportsByBuild for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CountTestReportsByBuild for %s returned err: %v", test.name, err)
			}

			if got != test.want {
				t.Errorf("CountTestReportsByBuild for %s is %d, want %d", test.name, got, test.want)
			}
		})
	}
}
