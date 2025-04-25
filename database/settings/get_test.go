// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestSettings_Engine_GetSettings(t *testing.T) {
	// setup types
	_settings := testSettings()
	_settings.SetID(1)
	_settings.SetCloneImage("target/vela-git-slim:latest")
	_settings.SetTemplateDepth(10)
	_settings.SetStarlarkExecLimit(100)
	_settings.SetRoutes([]string{"vela"})
	_settings.SetRepoAllowlist([]string{"octocat/hello-world"})
	_settings.SetScheduleAllowlist([]string{"*"})
	_settings.SetMaxDashboardRepos(10)
	_settings.SetCreatedAt(1)
	_settings.SetUpdatedAt(1)
	_settings.SetUpdatedBy("octocat")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.SettingsFromAPI(_settings)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "settings" WHERE id = $1 LIMIT $2`).WithArgs(1, 1).WillReturnRows(_rows)

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
		want     *settings.Platform
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _settings,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _settings,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetSettings(context.TODO())

			if test.failure {
				if err == nil {
					t.Errorf("GetSettings for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetSettings for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetSettings for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
