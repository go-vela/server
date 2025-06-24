// SPDX-License-Identifier: Apache-2.0

package testreport

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestTestReport_Engine_ListByRepo(t *testing.T) {
	// setup types
	_testReport := testutils.APITestReport()
	_testReport.SetID(1)
	_testReport.SetBuildID(1)
	_testReport.SetCreatedAt(1)

	_repo := testutils.APIRepo()
	_repo.SetID(1)

	_buildOne := testutils.APIBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepo(_repo)

	_postgres, _mock := testPostgres(t)
	ctx := context.TODO()
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query for the test_reports table
	_rows := sqlmock.NewRows([]string{"id", "repo_id", "build_id", "created_at"}).
		AddRow(1, 1, 1, 1)
	_mock.ExpectQuery(`SELECT testreports.* FROM "testreports" 
    JOIN builds ON testreports.build_id = builds.id 
    JOIN repos ON builds.repo_id = repos.id WHERE repo_id = $1 ORDER BY created_at DESC LIMIT $2`).
		WithArgs(1, 10).
		WillReturnRows(_rows)

	// Mock for Build preload query
	_buildRows := testutils.CreateMockRows([]any{*types.BuildFromAPI(_buildOne)})
	_mock.ExpectQuery(`SELECT * FROM "builds" WHERE "builds"."id" = $1`).
		WithArgs(1).
		WillReturnRows(_buildRows)

	// Mock for Repo preload query
	_repoRows := testutils.CreateMockRows([]any{*types.RepoFromAPI(_repo)})
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE "repos"."id" = $1`).
		WithArgs(1).
		WillReturnRows(_repoRows)

	// Mock for Owner preload query (using _repo.Owner)
	_ownerRows := testutils.CreateMockRows([]any{*types.UserFromAPI(_repo.GetOwner())})
	_mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."id" = $1`).
		WithArgs(1).
		WillReturnRows(_ownerRows)
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
	_, err = _sqlite.CreateTestReport(ctx, _testReport)
	if err != nil {
		t.Errorf("unable to create test report for sqlite: %v", err)
	}
	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     []*api.TestReport
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*api.TestReport{_testReport},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.TestReport{_testReport},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListTestReportsByRepo(ctx, _repo, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListTestReportsByRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListTestReportsByRepo for %s returned err: %v", test.name, err)
			}

			if len(got) != len(test.want) {
				t.Errorf("ListTestReportsByRepo for %s returned %d reports, want %d", test.name, len(got), len(test.want))
				return
			}

			if len(got) > 0 {
				// Check report fields
				if !reflect.DeepEqual(got[0].GetID(), test.want[0].GetID()) ||
					!reflect.DeepEqual(got[0].GetBuildID(), test.want[0].GetBuildID()) ||
					!reflect.DeepEqual(got[0].GetCreatedAt(), test.want[0].GetCreatedAt()) {
					t.Errorf("ListTestReportsByRepo for %s returned unexpected report values: got %v, want %v",
						test.name, got[0], test.want[0])
				}
			}
		})
	}
}
