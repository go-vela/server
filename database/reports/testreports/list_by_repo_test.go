// SPDX-License-Identifier: Apache-2.0

package testreports

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
)

func TestTestReports_Engine_ListByRepo(t *testing.T) {
	// setup types
	_testReport := testutils.APITestReport()
	_testReport.SetID(1)
	_testReport.SetBuildID(1)
	_testReport.SetCreated(1)

	_repo := testutils.APIRepo()
	_repo.SetID(1)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query for the test_reports table
	_rows := sqlmock.NewRows([]string{"id", "repo_id", "build_id", "created"}).
		AddRow(1, 1, 1, 1)
	_mock.ExpectQuery(`SELECT .* FROM "test_reports" WHERE repo_id = \$1 ORDER BY created DESC LIMIT \$2 OFFSET \$3`).
		WithArgs(1, 10, 0).
		WillReturnRows(_rows)

	// Mock for count query
	_countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	_mock.ExpectQuery(`SELECT count\(\*\) FROM "test_reports" WHERE repo_id = \$1`).
		WithArgs(1).
		WillReturnRows(_countRows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.Create(context.TODO(), _testReport)
	if err != nil {
		t.Errorf("unable to create test report for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure   bool
		name      string
		database  *Engine
		want      []*api.TestReport
		wantCount int64
	}{
		{
			failure:   false,
			name:      "postgres",
			database:  _postgres,
			want:      []*api.TestReport{_testReport},
			wantCount: 1,
		},
		{
			failure:   false,
			name:      "sqlite3",
			database:  _sqlite,
			want:      []*api.TestReport{_testReport},
			wantCount: 1,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, count, err := test.database.ListByRepo(context.TODO(), _repo, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListByRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListByRepo for %s returned err: %v", test.name, err)
			}

			if count != test.wantCount {
				t.Errorf("ListByRepo count for %s is %d, want %d", test.name, count, test.wantCount)
			}

			if len(got) != len(test.want) {
				t.Errorf("ListByRepo for %s returned %d reports, want %d", test.name, len(got), len(test.want))
				return
			}

			if len(got) > 0 {
				// Check report fields
				if !reflect.DeepEqual(got[0].GetID(), test.want[0].GetID()) ||
					!reflect.DeepEqual(got[0].GetBuildID(), test.want[0].GetBuildID()) ||
					!reflect.DeepEqual(got[0].GetCreated(), test.want[0].GetCreated()) {
					t.Errorf("ListByRepo for %s returned unexpected report values: got %v, want %v",
						test.name, got[0], test.want[0])
				}
			}
		})
	}
}
