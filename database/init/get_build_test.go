// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package init

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestInit_Engine_GetInitForRepo(t *testing.T) {
	// setup types
	_init := testInit()
	_init.SetID(1)
	_init.SetRepoID(1)
	_init.SetBuildID(1)
	_init.SetNumber(1)

	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "reporter", "name", "mimetype"}).
		AddRow(1, 1, 1, 1, "Foobar Runtime", "foobar", "text/plain")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "inits" WHERE build_id = $1 AND number = $2 LIMIT 1`).WithArgs(1, 1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateInit(_init)
	if err != nil {
		t.Errorf("unable to create test init for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     *library.Init
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _init,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _init,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetInitForBuild(_build, 1)

			if test.failure {
				if err == nil {
					t.Errorf("GetInitForBuild for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetInitForBuild for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetInitForBuild for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
