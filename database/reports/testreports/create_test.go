// SPDX-License-Identifier: Apache-2.0

package testreports

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestEngine_CreateTestReport(t *testing.T) {
	_testReport := testutils.APITestReport()
	_testReport.SetID(1)
	_testReport.SetBuildID(1)
	_testReport.SetCreated(1)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "testreports" ("build_id","created","id")
		VALUES ($1,$2,$3) RETURNING "id"`).
		WithArgs(1, 1, 1).
		WillReturnRows(_rows)

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
			got, err := test.database.Create(context.TODO(), _testReport)

			if test.failure {
				if err == nil {
					t.Errorf("Create for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("Create for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _testReport) {
				t.Errorf("Create for %s returned %v, want %v", test.name, got, _testReport)
			}
		})
	}
}
