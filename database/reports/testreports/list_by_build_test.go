// SPDX-License-Identifier: Apache-2.0

package testreports

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
	"testing"
)

func TestTestReports_Engine_ListByBuild(t *testing.T) {
	// setup types
	_testReport := testutils.APITestReport()
	_testReport.SetID(1)
	_testReport.SetBuildID(1)
	_testReport.SetCreated(1)

	_build := testutils.APIBuild()
	_build.SetID(1)

	_postgres, _mock := testPostgres(t)
	ctx := context.TODO()
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query for the test_reports table
	_rows := sqlmock.NewRows([]string{"id", "repo_id", "build_id", "created"}).
		AddRow(1, 1, 1, 1)
	_mock.ExpectQuery(`SELECT * FROM "testreports" WHERE build_id = $1 ORDER BY created DESC LIMIT $2`).
		WithArgs(1, 10).
		WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	// Create necessary SQLite tables for relationship testing
	err := _sqlite.client.AutoMigrate(&types.TestReport{}, &types.Build{})
	if err != nil {
		t.Errorf("unable to create tables for sqlite: %v", err)
	}

	// Set up build
	err = _sqlite.client.Table(constants.TableBuild).Create(types.BuildFromAPI(_build)).Error
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	// Then create the test report with the build_id
	_, err = _sqlite.Create(ctx, _testReport)
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
			got, err := test.database.ListByBuild(ctx, _build, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListByBuild for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListByBuild for %s returned err: %v", test.name, err)
			}

			if len(got) != len(test.want) {
				t.Errorf("ListByBuild for %s returned %d reports, want %d", test.name, len(got), len(test.want))
				return
			}

			if len(got) > 0 {
				// Check report fields
				if got[0].GetID() != test.want[0].GetID() ||
					got[0].GetBuildID() != test.want[0].GetBuildID() ||
					got[0].GetCreated() != test.want[0].GetCreated() {
					t.Errorf("ListByBuild for %s returned unexpected report values: got %v, want %v",
						test.name, got[0], test.want[0])
				}
			}
		})
	}
}
