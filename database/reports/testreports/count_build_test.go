// SPDX-License-Identifier: Apache-2.0

package testreports

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/server/database/testutils"
	"testing"
	"time"
)

func TestTestReports_Engine_CountByBuild(t *testing.T) {
	// setup types
	_testreport := testutils.APITestReport()
	_testreport.SetID(1)
	_testreport.SetBuildID(1)
	_testreport.SetCreated(1)

	_build := testutils.APIBuild()
	_build.SetID(1)

	_postgres, _mock := testPostgres(t)
	ctx := context.TODO()
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query for the test_reports table
	_mock.ExpectQuery(`SELECT count(*) FROM "testreports" WHERE build_id = $1 AND created < $2 AND created > $3`).
		WithArgs(1, AnyArgument{}, 0).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.Create(ctx, _testreport)
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
			got, err := test.database.CountByBuild(ctx, _build, filters, time.Now().Unix(), 0)
			if test.failure {
				if err == nil {
					t.Errorf("CountByBuild for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CountByBuild for %s returned err: %v", test.name, err)
			}

			if got != test.want {
				t.Errorf("CountByBuild for %s is %d, want %d", test.name, got, test.want)
			}
		})
	}

}
