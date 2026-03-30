// SPDX-License-Identifier: Apache-2.0

package installation

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestInstallation_Engine_GetInstallation(t *testing.T) {
	// setup types
	_installation := testutils.APIInstallation()
	_installation.SetInstallID(1)
	_installation.SetTarget("octocat")

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.InstallationFromAPI(_installation)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "installations" WHERE target = $1 LIMIT $2`).
		WithArgs("octocat", 1).
		WillReturnRows(_rows)

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateInstallation(context.TODO(), _installation)
	if err != nil {
		t.Errorf("unable to create test installation for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     *api.Installation
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _installation,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _installation,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetInstallation(context.TODO(), "octocat")

			if test.failure {
				if err == nil {
					t.Errorf("GetInstallation for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetInstallation for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetInstallation for %s: -want, +got: %s", test.name, diff)
			}
		})
	}
}
