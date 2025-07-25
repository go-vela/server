// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSettings_Engine_CreateSettings(t *testing.T) {
	// setup types
	_settings := testSettings()
	_settings.SetID(1)
	_settings.SetCloneImage("target/vela-git-slim:latest")
	_settings.SetTemplateDepth(10)
	_settings.SetStarlarkExecLimit(100)
	_settings.SetRoutes([]string{"vela"})
	_settings.SetRepoRoleMap(map[string]string{"admin": "admin", "triage": "read"})
	_settings.SetOrgRoleMap(map[string]string{"admin": "admin", "member": "read"})
	_settings.SetTeamRoleMap(map[string]string{"admin": "admin"})
	_settings.SetRepoAllowlist([]string{"octocat/hello-world"})
	_settings.SetScheduleAllowlist([]string{"*"})
	_settings.SetMaxDashboardRepos(10)
	_settings.SetQueueRestartLimit(30)
	_settings.SetCreatedAt(1)
	_settings.SetUpdatedAt(1)
	_settings.SetUpdatedBy("")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "settings" ("compiler","queue","scm","repo_allowlist","schedule_allowlist","max_dashboard_repos","queue_restart_limit","created_at","updated_at","updated_by","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`).
		WithArgs(`{"clone_image":{"String":"target/vela-git-slim:latest","Valid":true},"template_depth":{"Int64":10,"Valid":true},"starlark_exec_limit":{"Int64":100,"Valid":true}}`,
			`{"routes":["vela"]}`, `{"repo_role_map":{"admin":"admin","triage":"read"},"org_role_map":{"admin":"admin","member":"read"},"team_role_map":{"admin":"admin"}}`, `{"octocat/hello-world"}`, `{"*"}`, 10, 30, 1, 1, ``, 1).
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
			got, err := test.database.CreateSettings(context.TODO(), _settings)

			if test.failure {
				if err == nil {
					t.Errorf("CreateSettings for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateSettings for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _settings) {
				t.Errorf("CreateSettings for %s returned %s, want %s", test.name, got, _settings)
			}
		})
	}
}
