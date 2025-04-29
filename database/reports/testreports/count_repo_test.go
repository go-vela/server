// SPDX-License-Identifier: Apache-2.0

package testreports

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestTestReports_Engine_CountByRepo(t *testing.T) {
	// setup types
	_testReport := testutils.APITestReport()
	_testReport.SetID(1)
	_testReport.SetBuildID(1)
	_testReport.SetCreated(1)

	_repo := testutils.APIRepo()
	_repo.SetID(1)

	_buildOne := testutils.APIBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepo(_repo)

	_postgres, _mock := testPostgres(t)
	ctx := context.TODO()
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query for the test_reports table
	_mock.ExpectQuery(`SELECT count(*) FROM "testreports" 
    JOIN builds ON testreports.build_id = builds.id 
    JOIN repos ON builds.repo_id = repos.id WHERE repo_id = $1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	// Create necessary SQLite tables for relationship testing
	err := _sqlite.client.AutoMigrate(&types.TestReport{}, &types.Build{}, &types.Repo{}, &types.User{})
	if err != nil {
		t.Errorf("unable to create tables for sqlite: %v", err)
	}

	// Set up owner
	_owner := testutils.APIUser().Crop()
	_owner.SetID(1)
	err = _sqlite.client.Table(constants.TableUser).Create(types.UserFromAPI(_owner)).Error
	if err != nil {
		t.Errorf("unable to create test owner for sqlite: %v", err)
	}

	// Set up repo with owner
	_repo.SetOwner(_owner)
	err = _sqlite.client.Table(constants.TableRepo).Create(types.RepoFromAPI(_repo)).Error
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	// Set up build with repo
	_buildOne.SetRepo(_repo)
	err = _sqlite.client.Table(constants.TableBuild).Create(types.BuildFromAPI(_buildOne)).Error
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	// Then create the test report with the build_id
	_, err = _sqlite.Create(ctx, _testReport)
	if err != nil {
		t.Errorf("unable to create test report for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     int64
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     1,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     1,
		},
	}

	filters := map[string]interface{}{}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create test context
			ctx := context.Background()

			// Test the CountByRepo function
			got, err := test.database.CountByRepo(ctx, _repo, filters)

			// Check for expected errors
			if test.failure {
				if err == nil {
					t.Errorf("CountByRepo() error = nil, want error")
				}
				return
			}

			// Check for unexpected errors and validate results
			if err != nil {
				t.Errorf("CountByRepo() unexpected error: %v", err)
			}

			if got != test.want {
				t.Errorf("CountByRepo() got = %v, want %v", got, test.want)
			}
		})
	}
}
