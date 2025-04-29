// SPDX-License-Identifier: Apache-2.0

package testreports

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestTestReports_Engine_Update(t *testing.T) {
	// setup types
	_owner := testutils.APIUser()
	_owner.SetID(1)
	_owner.SetName("foo")
	_owner.SetToken("bar")

	_testReport := testutils.APITestReport()
	_testReport.SetID(1)
	_testReport.SetBuildID(1)
	_testReport.SetCreated(1)

	_postgres, _mock := testPostgres(t)
	ctx := context.TODO()
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query for the test_reports table
	_mock.ExpectExec(`UPDATE "testreports" SET "build_id"=$1,"created"=$2 WHERE "id" = $3`).
		WithArgs(1, 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

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
			got, err := test.database.Update(ctx, _testReport)

			if test.failure {
				if err == nil {
					t.Errorf("Update for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("Update for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got.GetID(), _testReport.GetID()) ||
				!reflect.DeepEqual(got.GetBuildID(), _testReport.GetBuildID()) ||
				!reflect.DeepEqual(got.GetCreated(), _testReport.GetCreated()) {
				t.Errorf("Update for %s returned unexpected report values: got %v, want %v", test.name, got, _testReport)
			}
		})
	}
}
