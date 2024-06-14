// SPDX-License-Identifier: Apache-2.0

package hook

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

func TestHook_Engine_ListHooksForRepo(t *testing.T) {
	// setup types
	_repo := testutils.APIRepo()
	_repo.SetID(1)
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")

	_build := testutils.APIBuild()
	_build.SetID(1)
	_build.SetRepo(_repo)

	_hookOne := testutils.APIHook()
	_hookOne.SetID(1)
	_hookOne.SetRepo(_repo)
	_hookOne.SetBuild(_build)
	_hookOne.SetNumber(1)
	_hookOne.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")
	_hookOne.SetWebhookID(1)

	_hookTwo := testutils.APIHook()
	_hookTwo.SetID(2)
	_hookTwo.SetRepo(_repo)
	_hookTwo.SetNumber(2)
	_hookTwo.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")
	_hookTwo.SetWebhookID(1)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "hooks" WHERE repo_id = $1`).WithArgs(1).WillReturnRows(_rows)

	// create expected result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "source_id", "created", "host", "event", "event_action", "branch", "error", "status", "link", "webhook_id"}).
		AddRow(2, 1, 2, 2, "c8da1302-07d6-11ea-882f-4893bca275b8", 0, "", "", "", "", "", "", "", 1).
		AddRow(1, 1, 1, 1, "c8da1302-07d6-11ea-882f-4893bca275b8", 0, "", "", "", "", "", "", "", 1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "hooks" WHERE repo_id = $1 ORDER BY id DESC LIMIT $2`).WithArgs(1, 10).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateHook(context.TODO(), _hookOne)
	if err != nil {
		t.Errorf("unable to create test hook for sqlite: %v", err)
	}

	_, err = _sqlite.CreateHook(context.TODO(), _hookTwo)
	if err != nil {
		t.Errorf("unable to create test hook for sqlite: %v", err)
	}

	err = _sqlite.client.AutoMigrate(&types.Repo{})
	if err != nil {
		t.Errorf("unable to create repo table for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableRepo).Create(types.RepoFromAPI(_repo)).Error
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	err = _sqlite.client.AutoMigrate(&types.Build{})
	if err != nil {
		t.Errorf("unable to create build table for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableBuild).Create(types.BuildFromAPI(_build)).Error
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     []*api.Hook
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*api.Hook{_hookTwo, _hookOne},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.Hook{_hookTwo, _hookOne},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _, err := test.database.ListHooksForRepo(context.TODO(), _repo, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListHooksForRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListHooksForRepo for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListHooksForRepo for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
