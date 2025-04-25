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
	_settings.SetRepoAllowlist([]string{"octocat/hello-world"})
	_settings.SetScheduleAllowlist([]string{"*"})
	_settings.SetMaxDashboardRepos(10)
	_settings.SetCreatedAt(1)
	_settings.SetUpdatedAt(1)
	_settings.SetUpdatedBy("octocat")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "settings" SET "compiler"=$1,"queue"=$2,"repo_allowlist"=$3,"schedule_allowlist"=$4,"max_dashboard_repos"=$5,"created_at"=$6,"updated_at"=$7,"updated_by"=$8 WHERE "id" = $9`).
		WithArgs(`{"clone_image":{"String":"target/vela-git-slim:latest","Valid":true},"template_depth":{"Int64":10,"Valid":true},"starlark_exec_limit":{"Int64":100,"Valid":true}}`,
			`{"routes":["vela","large"]}`, `{"octocat/hello-world"}`, `{"*"}`, 10, 1, testutils.AnyArgument{}, "octocat", 1).
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
