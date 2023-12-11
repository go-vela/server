// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

func TestRepo_Engine_ListReposForOrg(t *testing.T) {
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

	_repoOne := testRepo()
	_repoOne.SetID(1)
	_repoOne.SetUserID(1)
	_repoOne.SetHash("baz")
	_repoOne.SetOrg("foo")
	_repoOne.SetName("bar")
	_repoOne.SetFullName("foo/bar")
	_repoOne.SetVisibility("public")
	_repoOne.SetPipelineType("yaml")
	_repoOne.SetTopics([]string{})
	_repoOne.SetAllowEvents(library.NewEventsFromMask(1))

	_repoTwo := testRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetUserID(1)
	_repoTwo.SetHash("bar")
	_repoTwo.SetOrg("foo")
	_repoTwo.SetName("baz")
	_repoTwo.SetFullName("foo/baz")
	_repoTwo.SetVisibility("public")
	_repoTwo.SetPipelineType("yaml")
	_repoTwo.SetTopics([]string{})
	_repoTwo.SetAllowEvents(library.NewEventsFromMask(1))

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected name count query result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the name count query
	_mock.ExpectQuery(`SELECT count(*) FROM "repos" WHERE org = $1`).WithArgs("foo").WillReturnRows(_rows)

	// create expected name query result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "user_id", "hash", "org", "name", "full_name", "link", "clone", "branch", "topics", "timeout", "counter", "visibility", "private", "trusted", "active", "allow_pull", "allow_push", "allow_deploy", "allow_tag", "allow_comment", "allow_events", "pipeline_type", "previous_name", "approve_build"}).
		AddRow(1, 1, "baz", "foo", "bar", "foo/bar", "", "", "", "{}", 0, 0, "public", false, false, false, false, false, false, false, false, 1, "yaml", nil, nil).
		AddRow(2, 1, "bar", "foo", "baz", "foo/baz", "", "", "", "{}", 0, 0, "public", false, false, false, false, false, false, false, false, 1, "yaml", nil, nil)

	// ensure the mock expects the name query
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE org = $1 ORDER BY name LIMIT 10`).WithArgs("foo").WillReturnRows(_rows)

	// create expected latest count query result in mock
	_rows = sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the latest count query
	_mock.ExpectQuery(`SELECT count(*) FROM "repos" WHERE org = $1`).WithArgs("foo").WillReturnRows(_rows)

	// create expected latest query result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "user_id", "hash", "org", "name", "full_name", "link", "clone", "branch", "topics", "timeout", "counter", "visibility", "private", "trusted", "active", "allow_pull", "allow_push", "allow_deploy", "allow_tag", "allow_comment", "allow_events", "pipeline_type", "previous_name", "approve_build"}).
		AddRow(1, 1, "baz", "foo", "bar", "foo/bar", "", "", "", "{}", 0, 0, "public", false, false, false, false, false, false, false, false, 1, "yaml", nil, nil).
		AddRow(2, 1, "bar", "foo", "baz", "foo/baz", "", "", "", "{}", 0, 0, "public", false, false, false, false, false, false, false, false, 1, "yaml", nil, nil)

	// ensure the mock expects the latest query
	_mock.ExpectQuery(`SELECT repos.* FROM "repos" LEFT JOIN (SELECT repos.id, MAX(builds.created) AS latest_build FROM "builds" INNER JOIN repos repos ON builds.repo_id = repos.id WHERE repos.org = $1 GROUP BY "repos"."id") t on repos.id = t.id ORDER BY latest_build DESC NULLS LAST LIMIT 10`).WithArgs("foo").WillReturnRows(_rows)

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

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		sort     string
		database *engine
		want     []*library.Repo
	}{
		{
			failure:  false,
			name:     "postgres with name",
			database: _postgres,
			sort:     "name",
			want:     []*library.Repo{_repoOne, _repoTwo},
		},
		{
			failure:  false,
			name:     "postgres with latest",
			database: _postgres,
			sort:     "latest",
			want:     []*library.Repo{_repoOne, _repoTwo},
		},
		{
			failure:  false,
			name:     "sqlite with name",
			database: _sqlite,
			sort:     "name",
			want:     []*library.Repo{_repoOne, _repoTwo},
		},
		{
			failure:  false,
			name:     "sqlite with latest",
			database: _sqlite,
			sort:     "latest",
			want:     []*library.Repo{_repoOne, _repoTwo},
		},
	}

	filters := map[string]interface{}{}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _, err := test.database.ListReposForOrg(context.TODO(), "foo", test.sort, filters, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListReposForOrg for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListReposForOrg for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListReposForOrg for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
