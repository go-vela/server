// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestSettings_Engine_UpdateSettings(t *testing.T) {
	// setup types
	_settings := testSettings()
	_settings.SetID(1)
	_settings.SetCloneImage("target/vela-git-slim:latest")
	_settings.SetTemplateDepth(10)
	_settings.SetStarlarkExecLimit(100)
	_settings.SetRoutes([]string{"vela", "large"})
	_settings.SetRepoRoleMap(map[string]string{"admin": "admin", "triage": "read"})
	_settings.SetOrgRoleMap(map[string]string{"admin": "admin", "member": "read"})
	_settings.SetTeamRoleMap(map[string]string{"admin": "admin"})
	_settings.SetRepoAllowlist([]string{"octocat/hello-world"})
	_settings.SetScheduleAllowlist([]string{"*"})
	_settings.SetMaxDashboardRepos(10)
	_settings.SetQueueRestartLimit(30)
	_settings.SetCreatedAt(1)
	_settings.SetUpdatedAt(1)
	_settings.SetUpdatedBy("octocat")

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "settings" SET "compiler"=$1,"queue"=$2,"scm"=$3,"repo_allowlist"=$4,"schedule_allowlist"=$5,"max_dashboard_repos"=$6,"queue_restart_limit"=$7,"created_at"=$8,"updated_at"=$9,"updated_by"=$10 WHERE "id" = $11`).
		WithArgs(`{"clone_image":{"String":"target/vela-git-slim:latest","Valid":true},"template_depth":{"Int64":10,"Valid":true},"starlark_exec_limit":{"Int64":100,"Valid":true}}`,
			`{"routes":["vela","large"]}`, `{"repo_role_map":{"admin":"admin","triage":"read"},"org_role_map":{"admin":"admin","member":"read"},"team_role_map":{"admin":"admin"}}`, `{"octocat/hello-world"}`, `{"*"}`, 10, 30, 1, testutils.AnyArgument{}, "octocat", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateSettings(context.TODO(), _settings)
	if err != nil {
		t.Errorf("unable to create test settings for sqlite: %v", err)
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
			got, err := test.database.UpdateSettings(context.TODO(), _settings)
			got.SetUpdatedAt(_settings.GetUpdatedAt())

			if test.failure {
				if err == nil {
					t.Errorf("UpdateSettings for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateSettings for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _settings) {
				t.Errorf("UpdateSettings for %s returned %s, want %s", test.name, got, _settings)
			}
		})
	}
}
