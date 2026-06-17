// SPDX-License-Identifier: Apache-2.0

package limits

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestLimits_Engine_CreateOrgBuildLimit(t *testing.T) {
	// setup types
	_orgBuildLimit := testOrgBuildLimit()

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "org_build_limits" ("org","build_limit","created_at","updated_at","updated_by","id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`).
		WithArgs("github", 30, 1, 1, "octocat", 1).
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
			got, err := test.database.CreateOrgBuildLimit(context.TODO(), _orgBuildLimit)

			if test.failure {
				if err == nil {
					t.Errorf("CreateOrgBuildLimit for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateOrgBuildLimit for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _orgBuildLimit) {
				t.Errorf("CreateOrgBuildLimit for %s returned %s, want %s", test.name, got, _orgBuildLimit)
			}
		})
	}
}
