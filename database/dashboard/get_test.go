// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestRepo_Engine_GetDashboard(t *testing.T) {
	// setup types
	_dashRepo := testDashboardRepo()
	_dashRepo.SetID(1)
	_dashRepos := []*library.DashboardRepo{_dashRepo}

	_dashboard := testDashboard()
	_dashboard.SetID("abc-123")
	_dashboard.SetName("dash")
	_dashboard.SetRepos(_dashRepos)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "name", "created_at", "created_by", "updated_at", "updated_by", "admins", "repos"},
	).AddRow(1, "dash", 1, "user", 1, "user", "{}", `[{"id":1}]`)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "dashboards" WHERE id = $1 LIMIT 1`).WithArgs(1).WillReturnRows(_rows)

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
		database *engine
		want     *library.Dashboard
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
			got, err := test.database.GetDashboard(context.TODO(), "123-abc")

			if test.failure {
				if err == nil {
					t.Errorf("GetDashboard for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetDashboard for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetDashboard for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
