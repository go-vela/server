// SPDX-License-Identifier: Apache-2.0

package database

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/server/database/build"
	"github.com/go-vela/server/database/executable"
	"github.com/go-vela/server/database/hook"
	"github.com/go-vela/server/database/log"
	"github.com/go-vela/server/database/pipeline"
	"github.com/go-vela/server/database/repo"
	"github.com/go-vela/server/database/schedule"
	"github.com/go-vela/server/database/secret"
	"github.com/go-vela/server/database/service"
	"github.com/go-vela/server/database/step"
	"github.com/go-vela/server/database/user"
	"github.com/go-vela/server/database/worker"
)

func TestDatabase_Engine_NewResources(t *testing.T) {
	_postgres, _mock := testPostgres(t)
	defer _postgres.Close()

	// ensure the mock expects the build queries
	_mock.ExpectExec(build.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(build.CreateCreatedIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(build.CreateRepoIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(build.CreateSourceIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(build.CreateStatusIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(build.CreateSenderIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	// ensure the mock expects the build executable queries
	_mock.ExpectExec(executable.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	// ensure the mock expects the hook queries
	_mock.ExpectExec(hook.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(hook.CreateRepoIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	// ensure the mock expects the log queries
	_mock.ExpectExec(log.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(log.CreateBuildIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	// ensure the mock expects the pipeline queries
	_mock.ExpectExec(pipeline.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(pipeline.CreateRepoIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	// ensure the mock expects the repo queries
	_mock.ExpectExec(repo.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(repo.CreateOrgNameIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	// ensure the mock expects the schedule queries
	_mock.ExpectExec(schedule.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(schedule.CreateRepoIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	// ensure the mock expects the secret queries
	_mock.ExpectExec(secret.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(secret.CreateTypeOrgRepo).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(secret.CreateTypeOrgTeam).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(secret.CreateTypeOrg).WillReturnResult(sqlmock.NewResult(1, 1))
	// ensure the mock expects the service queries
	_mock.ExpectExec(service.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	// ensure the mock expects the step queries
	_mock.ExpectExec(step.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	// ensure the mock expects the user queries
	_mock.ExpectExec(user.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(user.CreateUserRefreshIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	// ensure the mock expects the worker queries
	_mock.ExpectExec(worker.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(worker.CreateHostnameAddressIndex).WillReturnResult(sqlmock.NewResult(1, 1))

	// create a test database without mocking the call
	_unmocked, _ := testPostgres(t)

	_sqlite := testSqlite(t)
	defer _sqlite.Close()

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
	}{
		{
			name:     "success with postgres",
			failure:  false,
			database: _postgres,
		},
		{
			name:     "success with sqlite3",
			failure:  false,
			database: _sqlite,
		},
		{
			name:     "failure without mocked call",
			failure:  true,
			database: _unmocked,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.database.NewResources(context.TODO())

			if test.failure {
				if err == nil {
					t.Errorf("NewResources for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("NewResources for %s returned err: %v", test.name, err)
			}
		})
	}
}
