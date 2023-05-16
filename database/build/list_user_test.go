// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

func TestBuild_Engine_ListBuildsForUser(t *testing.T) {
	// setup types
	_buildOne := new(library.Build)
	_buildOne.SetID(1)
	_buildOne.SetBuildID(1)
	_buildOne.SetNumber(1)
	_buildOne.SetCreated(time.Now().UTC().Unix())

	_buildTwo := new(library.Build)
	_buildTwo.SetID(2)
	_buildTwo.SetBuildID(2)
	_buildTwo.SetNumber(1)
	_buildTwo.SetCreated(time.Now().UTC().Unix())

	_buildOne := testBuild()
	_buildOne.SetID(1)
	_buildOne.SetUserID(1)
	_buildOne.SetHash("baz")
	_buildOne.SetOrg("foo")
	_buildOne.SetName("bar")
	_buildOne.SetFullName("foo/bar")
	_buildOne.SetVisibility("public")
	_buildOne.SetPipelineType("yaml")
	_buildOne.SetTopics([]string{})

	_buildTwo := testBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetUserID(1)
	_buildTwo.SetHash("baz")
	_buildTwo.SetOrg("bar")
	_buildTwo.SetName("foo")
	_buildTwo.SetFullName("bar/foo")
	_buildTwo.SetVisibility("public")
	_buildTwo.SetPipelineType("yaml")
	_buildTwo.SetTopics([]string{})

	_user := new(library.User)
	_user.SetID(1)
	_user.SetName("foo")
	_user.SetToken("bar")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected name count query result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the name count query
	_mock.ExpectQuery(`SELECT count(*) FROM "builds" WHERE user_id = $1`).WithArgs(1).WillReturnRows(_rows)

	// create expected name query result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "user_id", "hash", "org", "name", "full_name", "link", "clone", "branch", "topics", "timeout", "counter", "visibility", "private", "trusted", "active", "allow_pull", "allow_push", "allow_deploy", "allow_tag", "allow_comment", "pipeline_type", "previous_name"}).
		AddRow(1, 1, "baz", "foo", "bar", "foo/bar", "", "", "", "{}", 0, 0, "public", false, false, false, false, false, false, false, false, "yaml", nil).
		AddRow(2, 1, "baz", "bar", "foo", "bar/foo", "", "", "", "{}", 0, 0, "public", false, false, false, false, false, false, false, false, "yaml", nil)

	// ensure the mock expects the name query
	_mock.ExpectQuery(`SELECT * FROM "builds" WHERE user_id = $1 ORDER BY name LIMIT 10`).WithArgs(1).WillReturnRows(_rows)

	// create expected latest count query result in mock
	_rows = sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the latest count query
	_mock.ExpectQuery(`SELECT count(*) FROM "builds" WHERE user_id = $1`).WithArgs(1).WillReturnRows(_rows)

	// create expected latest query result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "user_id", "hash", "org", "name", "full_name", "link", "clone", "branch", "topics", "timeout", "counter", "visibility", "private", "trusted", "active", "allow_pull", "allow_push", "allow_deploy", "allow_tag", "allow_comment", "pipeline_type", "previous_name"}).
		AddRow(1, 1, "baz", "foo", "bar", "foo/bar", "", "", "", "{}", 0, 0, "public", false, false, false, false, false, false, false, false, "yaml", nil).
		AddRow(2, 1, "baz", "bar", "foo", "bar/foo", "", "", "", "{}", 0, 0, "public", false, false, false, false, false, false, false, false, "yaml", nil)

	// ensure the mock expects the latest query
	_mock.ExpectQuery(`SELECT builds.* FROM "builds" LEFT JOIN (SELECT builds.id, MAX(builds.created) AS latest_build FROM "builds" INNER JOIN builds builds ON builds.repo_id = builds.id WHERE builds.user_id = $1 GROUP BY "builds"."id") t on builds.id = t.id ORDER BY latest_build DESC NULLS LAST LIMIT 10`).WithArgs(1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateBuild(_buildOne)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	err = _sqlite.CreateBuild(_buildTwo)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
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
		want     []*library.Build
	}{
		{
			failure:  false,
			name:     "postgres with name",
			database: _postgres,
			sort:     "name",
			want:     []*library.Build{_buildOne, _buildTwo},
		},
		{
			failure:  false,
			name:     "postgres with latest",
			database: _postgres,
			sort:     "latest",
			want:     []*library.Build{_buildOne, _buildTwo},
		},
		{
			failure:  false,
			name:     "sqlite with name",
			database: _sqlite,
			sort:     "name",
			want:     []*library.Build{_buildOne, _buildTwo},
		},
		{
			failure:  false,
			name:     "sqlite with latest",
			database: _sqlite,
			sort:     "latest",
			want:     []*library.Build{_buildOne, _buildTwo},
		},
	}

	filters := map[string]interface{}{}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _, err := test.database.ListBuildsForUser(_user, test.sort, filters, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListBuildsForUser for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListBuildsForUser for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListBuildsForUser for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
