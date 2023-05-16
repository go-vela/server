// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestBuild_Engine_ListBuilds(t *testing.T) {
	// setup types
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

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "builds"`).WillReturnRows(_rows)

	// create expected result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "user_id", "hash", "org", "name", "full_name", "link", "clone", "branch", "topics", "timeout", "counter", "visibility", "private", "trusted", "active", "allow_pull", "allow_push", "allow_deploy", "allow_tag", "allow_comment", "pipeline_type", "previous_name"}).
		AddRow(1, 1, "baz", "foo", "bar", "foo/bar", "", "", "", "{}", 0, 0, "public", false, false, false, false, false, false, false, false, "yaml", nil).
		AddRow(2, 1, "baz", "bar", "foo", "bar/foo", "", "", "", "{}", 0, 0, "public", false, false, false, false, false, false, false, false, "yaml", nil)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "builds"`).WillReturnRows(_rows)

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

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     []*library.Build
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*library.Build{_buildOne, _buildTwo},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*library.Build{_buildOne, _buildTwo},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListBuilds()

			if test.failure {
				if err == nil {
					t.Errorf("ListBuilds for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListBuilds for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListBuilds for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
