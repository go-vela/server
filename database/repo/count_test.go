// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestRepo_Engine_CountRepos(t *testing.T) {
	// setup types
	_repoOne := testutils.APIRepo()
	_repoOne.SetID(1)
	_repoOne.GetOwner().SetID(1)
	_repoOne.SetHash("baz")
	_repoOne.SetOrg("foo")
	_repoOne.SetName("bar")
	_repoOne.SetFullName("foo/bar")
	_repoOne.SetVisibility("public")

	_repoTwo := testutils.APIRepo()
	_repoTwo.SetID(2)
	_repoTwo.GetOwner().SetID(1)
	_repoTwo.SetHash("baz")
	_repoTwo.SetOrg("bar")
	_repoTwo.SetName("foo")
	_repoTwo.SetFullName("bar/foo")
	_repoTwo.SetVisibility("public")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "repos"`).WillReturnRows(_rows)

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
			want:     2,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     2,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.CountRepos(context.TODO())

			if test.failure {
				if err == nil {
					t.Errorf("CountRepos for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CountRepos for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("CountRepos for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
