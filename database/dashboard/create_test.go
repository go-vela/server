// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestDashboard_Engine_CreateDashboard(t *testing.T) {
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

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow("c8da1302-07d6-11ea-882f-4893bca275b8")

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "dashboards"
("name","created_at","created_by","updated_at","updated_by","admins","repos","id")
VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`).
		WithArgs("dash", 1, "user1", 1, "user2", "{}", `[{"id":1,"branches":["main"],"events":["push"]}]`, "c8da1302-07d6-11ea-882f-4893bca275b8").
		WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

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
			got, err := test.database.CreateDashboard(context.TODO(), _dashboard)

			if test.failure {
				if err == nil {
					t.Errorf("CreateDashboard for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateDashboard for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _dashboard) {
				t.Errorf("CreateDashboard for %s returned %s, want %s", test.name, got, _dashboard)
			}
		})
	}
}
