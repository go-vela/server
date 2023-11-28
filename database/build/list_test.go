// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestBuild_Engine_ListBuilds(t *testing.T) {
	// setup types
	_buildOne := testBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepoID(1)
	_buildOne.SetNumber(1)
	_buildOne.SetDeployPayload(nil)

	_buildTwo := testBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepoID(1)
	_buildTwo.SetNumber(2)
	_buildTwo.SetDeployPayload(nil)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "builds"`).WillReturnRows(_rows)

	// create expected result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "repo_id", "pipeline_id", "number", "parent", "event", "event_action", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "approved_at", "approved_by", "timestamp"}).
		AddRow(1, 1, nil, 1, 0, "", "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0, "", 0).
		AddRow(2, 1, nil, 2, 0, "", "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0, "", 0)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "builds"`).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateBuild(context.TODO(), _buildOne)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	_, err = _sqlite.CreateBuild(context.TODO(), _buildTwo)
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
			got, err := test.database.ListBuilds(context.TODO())

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
