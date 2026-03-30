// SPDX-License-Identifier: Apache-2.0

package installation

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestInstallation_Engine_CreateInstallation(t *testing.T) {
	// setup types
	_installation := testutils.APIInstallation()
	_installation.SetInstallID(1)
	_installation.SetTarget("octocat")

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`INSERT INTO "installations" ("install_id","target") VALUES ($1,$2)`).
		WithArgs(1, "octocat").
		WillReturnResult(sqlmock.NewResult(1, 1))

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
			_, err := test.database.CreateInstallation(context.TODO(), _installation)

			if test.failure {
				if err == nil {
					t.Errorf("CreateInstallation for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateInstallation for %s returned err: %v", test.name, err)
			}
		})
	}
}
