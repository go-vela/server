// SPDX-License-Identifier: Apache-2.0

package limits

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestLimits_Engine_UpdateOrgBuildLimit(t *testing.T) {
	// setup types
	_orgBuildLimit := testOrgBuildLimit()

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "org_build_limits" SET "org"=$1,"build_limit"=$2,"created_at"=$3,"updated_at"=$4,"updated_by"=$5 WHERE "id" = $6`).
		WithArgs("github", 30, 1, testutils.AnyArgument{}, "octocat", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateOrgBuildLimit(context.TODO(), _orgBuildLimit)
	if err != nil {
		t.Errorf("unable to create test org build limit for sqlite: %v", err)
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
			got, err := test.database.UpdateOrgBuildLimit(context.TODO(), _orgBuildLimit)
			got.SetUpdatedAt(_orgBuildLimit.GetUpdatedAt())

			if test.failure {
				if err == nil {
					t.Errorf("UpdateOrgBuildLimit for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateOrgBuildLimit for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _orgBuildLimit) {
				t.Errorf("UpdateOrgBuildLimit for %s returned %s, want %s", test.name, got, _orgBuildLimit)
			}
		})
	}
}
