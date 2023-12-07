// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestDashboard_Engine_DeleteDashboard(t *testing.T) {
	// setup types
	_dashboard := testDashboard()
	_dashboard.SetID("c8da1302-07d6-11ea-882f-4893bca275b8")
	_dashboard.SetName("vela")
	_dashboard.SetCreatedAt(1)
	_dashboard.SetCreatedBy("user1")
	_dashboard.SetUpdatedAt(1)
	_dashboard.SetUpdatedBy("user2")
	_dashboard.SetAdmins([]string{"octocat"})
	_dashboard.SetRepos([]*library.DashboardRepo{testDashboardRepo()})

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`DELETE FROM "dashboards" WHERE "dashboards"."id" = $1`).
		WithArgs("c8da1302-07d6-11ea-882f-4893bca275b8").
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
			err = test.database.DeleteDashboard(context.TODO(), _dashboard)

			if test.failure {
				if err == nil {
					t.Errorf("DeleteDashboard for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("DeleteDashboard for %s returned err: %v", test.name, err)
			}
		})
	}
}
