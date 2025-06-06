// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestRepo_Engine_GetDashboard(t *testing.T) {
	// setup types
	_dashRepo := new(api.DashboardRepo)
	_dashRepo.SetID(1)
	_dashRepo.SetBranches([]string{"main"})
	_dashRepo.SetEvents([]string{"push"})
	_dashRepos := []*api.DashboardRepo{_dashRepo}

	_admin := new(api.User)
	_admin.SetID(1)
	_admin.SetName("octocat")
	_admin.SetActive(true)
	_admins := []*api.User{_admin}

	_dashboard := testutils.APIDashboard()
	_dashboard.SetID("c8da1302-07d6-11ea-882f-4893bca275b8")
	_dashboard.SetName("dash")
	_dashboard.SetCreatedAt(1)
	_dashboard.SetCreatedBy("user1")
	_dashboard.SetUpdatedAt(1)
	_dashboard.SetUpdatedBy("user2")
	_dashboard.SetRepos(_dashRepos)
	_dashboard.SetAdmins(_admins)

	// uuid, _ := uuid.Parse("c8da1302-07d6-11ea-882f-4893bca275b8")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.DashboardFromAPI(_dashboard)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "dashboards" WHERE id = $1 LIMIT $2`).WithArgs("c8da1302-07d6-11ea-882f-4893bca275b8", 1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateDashboard(context.TODO(), _dashboard)
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     *api.Dashboard
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _dashboard,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _dashboard,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetDashboard(context.TODO(), "c8da1302-07d6-11ea-882f-4893bca275b8")

			if test.failure {
				if err == nil {
					t.Errorf("GetDashboard for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetDashboard for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetDashboard mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
