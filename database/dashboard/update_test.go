// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
	"github.com/google/go-cmp/cmp"
)

func TestDashboard_Engine_UpdateDashboard(t *testing.T) {
	// setup types
	_dashRepo := new(library.DashboardRepo)
	_dashRepo.SetID(1)
	_dashRepo.SetBranches([]string{"main"})
	_dashRepo.SetEvents([]string{"push"})
	_dashRepos := []*library.DashboardRepo{_dashRepo}

	_dashboard := testDashboard()
	_dashboard.SetID("c8da1302-07d6-11ea-882f-4893bca275b8")
	_dashboard.SetName("dash")
	_dashboard.SetCreatedAt(1)
	_dashboard.SetCreatedBy("user1")
	_dashboard.SetUpdatedAt(1)
	_dashboard.SetUpdatedBy("user2")
	_dashboard.SetAdmins([]string{})
	_dashboard.SetRepos(_dashRepos)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "dashboards"
SET "name"=$1,"created_at"=$2,"created_by"=$3,"updated_at"=$4,"updated_by"=$5,"admins"=$6,"repos"=$7 WHERE "id" = $8`).
		WithArgs("dash", 1, "user1", NowTimestamp{}, "user2", "{}", `[{"id":1,"branches":["main"],"events":["push"]}]`, "c8da1302-07d6-11ea-882f-4893bca275b8").
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateDashboard(context.TODO(), _dashboard)
	if err != nil {
		t.Errorf("unable to create test dashboard for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
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
			got, err := test.database.UpdateDashboard(context.TODO(), _dashboard)
			_dashboard.SetUpdatedAt(got.GetUpdatedAt())

			if test.failure {
				if err == nil {
					t.Errorf("UpdateDashboard for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateDashboard for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(got, _dashboard); diff != "" {
				t.Errorf("GetDashboard mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
