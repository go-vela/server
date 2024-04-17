// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/user"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

func TestRepo_Engine_ListReposForUser(t *testing.T) {
	// setup types
	_buildOne := new(library.Build)
	_buildOne.SetID(1)
	_buildOne.SetRepoID(1)
	_buildOne.SetNumber(1)
	_buildOne.SetCreated(time.Now().UTC().Unix())

	_buildTwo := new(library.Build)
	_buildTwo.SetID(2)
	_buildTwo.SetRepoID(2)
	_buildTwo.SetNumber(1)
	_buildTwo.SetCreated(time.Now().UTC().Unix())

	_repoOne := testAPIRepo()
	_repoOne.SetID(1)
	_repoOne.GetOwner().SetID(1)
	_repoOne.SetHash("baz")
	_repoOne.SetOrg("foo")
	_repoOne.SetName("bar")
	_repoOne.SetFullName("foo/bar")
	_repoOne.SetVisibility("public")
	_repoOne.SetPipelineType("yaml")
	_repoOne.SetTopics([]string{})
	_repoOne.SetAllowEvents(api.NewEventsFromMask(1))

	_repoTwo := testAPIRepo()
	_repoTwo.SetID(2)
	_repoTwo.GetOwner().SetID(1)
	_repoTwo.SetHash("baz")
	_repoTwo.SetOrg("bar")
	_repoTwo.SetName("foo")
	_repoTwo.SetFullName("bar/foo")
	_repoTwo.SetVisibility("public")
	_repoTwo.SetPipelineType("yaml")
	_repoTwo.SetTopics([]string{})
	_repoTwo.SetAllowEvents(api.NewEventsFromMask(1))

	_owner := testOwner()
	_owner.SetID(1)
	_owner.SetName("foo")

	_repoOne.SetOwner(_owner)
	_repoTwo.SetOwner(_owner)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected name count query result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the name count query
	_mock.ExpectQuery(`SELECT count(*) FROM "repos" WHERE user_id = $1`).WithArgs(1).WillReturnRows(_rows)

	// create expected name query result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "user_id", "hash", "org", "name", "full_name", "link", "clone", "branch", "topics", "timeout", "counter", "visibility", "private", "trusted", "active", "allow_events", "pipeline_type", "previous_name", "approve_build"}).
		AddRow(1, 1, "baz", "foo", "bar", "foo/bar", "", "", "", "{}", 0, 0, "public", false, false, false, 1, "yaml", nil, nil).
		AddRow(2, 1, "baz", "bar", "foo", "bar/foo", "", "", "", "{}", 0, 0, "public", false, false, false, 1, "yaml", nil, nil)

	_userRows := sqlmock.NewRows(
		[]string{"id", "name", "token", "hash", "active", "admin"}).
		AddRow(1, "foo", "bar", "baz", false, false)

	// ensure the mock expects the name query
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE user_id = $1 ORDER BY name LIMIT $2`).WithArgs(1, 10).WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."id" = $1`).WithArgs(1).WillReturnRows(_userRows)

	// create expected latest count query result in mock
	_rows = sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the latest count query
	_mock.ExpectQuery(`SELECT count(*) FROM "repos" WHERE user_id = $1`).WithArgs(1).WillReturnRows(_rows)

	// create expected latest query result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "user_id", "hash", "org", "name", "full_name", "link", "clone", "branch", "topics", "timeout", "counter", "visibility", "private", "trusted", "active", "allow_events", "pipeline_type", "previous_name", "approve_build"}).
		AddRow(1, 1, "baz", "foo", "bar", "foo/bar", "", "", "", "{}", 0, 0, "public", false, false, false, 1, "yaml", nil, nil).
		AddRow(2, 1, "baz", "bar", "foo", "bar/foo", "", "", "", "{}", 0, 0, "public", false, false, false, 1, "yaml", nil, nil)

	_userRows = sqlmock.NewRows(
		[]string{"id", "name", "token", "hash", "active", "admin"}).
		AddRow(1, "foo", "bar", "baz", false, false)

	// ensure the mock expects the latest query
	_mock.ExpectQuery(`SELECT repos.* FROM "repos" LEFT JOIN (SELECT repos.id, MAX(builds.created) AS latest_build FROM "builds" INNER JOIN repos repos ON builds.repo_id = repos.id WHERE repos.user_id = $1 GROUP BY "repos"."id") t on repos.id = t.id ORDER BY latest_build DESC NULLS LAST LIMIT $2`).WithArgs(1, 10).WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."id" = $1`).WithArgs(1).WillReturnRows(_userRows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateRepo(context.TODO(), _repoOne)
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	_, err = _sqlite.CreateRepo(context.TODO(), _repoTwo)
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	err = _sqlite.client.AutoMigrate(&database.Build{})
	if err != nil {
		t.Errorf("unable to create build table for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableBuild).Create(database.BuildFromLibrary(_buildOne).Crop()).Error
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableBuild).Create(database.BuildFromLibrary(_buildTwo).Crop()).Error
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	err = _sqlite.client.AutoMigrate(&user.User{})
	if err != nil {
		t.Errorf("unable to create build table for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableUser).Create(user.FromAPI(_owner)).Error
	if err != nil {
		t.Errorf("unable to create test user for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		sort     string
		database *engine
		want     []*api.Repo
	}{
		{
			failure:  false,
			name:     "postgres with name",
			database: _postgres,
			sort:     "name",
			want:     []*api.Repo{_repoOne, _repoTwo},
		},
		{
			failure:  false,
			name:     "postgres with latest",
			database: _postgres,
			sort:     "latest",
			want:     []*api.Repo{_repoOne, _repoTwo},
		},
		{
			failure:  false,
			name:     "sqlite with name",
			database: _sqlite,
			sort:     "name",
			want:     []*api.Repo{_repoOne, _repoTwo},
		},
		{
			failure:  false,
			name:     "sqlite with latest",
			database: _sqlite,
			sort:     "latest",
			want:     []*api.Repo{_repoOne, _repoTwo},
		},
	}

	filters := map[string]interface{}{}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _, err := test.database.ListReposForUser(context.TODO(), _owner, test.sort, filters, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListReposForUser for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListReposForUser for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListReposForUser for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
