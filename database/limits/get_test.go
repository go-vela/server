// SPDX-License-Identifier: Apache-2.0

package limits

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestLimits_Engine_GetOrgBuildLimit(t *testing.T) {
	// setup types
	_orgBuildLimit := testOrgBuildLimit()

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.OrgBuildLimitFromAPI(_orgBuildLimit)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "org_build_limits" WHERE org = $1 LIMIT $2`).
		WithArgs("github", 1).
		WillReturnRows(_rows)

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
		want     *api.OrgBuildLimit
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _orgBuildLimit,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _orgBuildLimit,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetOrgBuildLimit(context.TODO(), "github")

			if test.failure {
				if err == nil {
					t.Errorf("GetOrgBuildLimit for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetOrgBuildLimit for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetOrgBuildLimit mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
