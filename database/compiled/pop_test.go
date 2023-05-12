// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package compiled

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestCompiled_Engine_PopCompiled(t *testing.T) {
	// setup types
	_compiled := testCompiled()
	_compiled.SetID(1)
	_compiled.SetBuildID(1)
	_compiled.SetData([]byte("foo"))

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "build_id", "data"}).
		AddRow(1, 1, []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69})

	// ensure the mock expects the query
	_mock.ExpectQuery(`DELETE FROM "compiled" WHERE build_id = $1 RETURNING *`).WithArgs(1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateCompiled(_compiled)
	if err != nil {
		t.Errorf("unable to create test pipeline for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     *library.Compiled
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _compiled,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _compiled,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.PopCompiled(1)

			if test.failure {
				if err == nil {
					t.Errorf("GetPipeline for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetPipeline for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetPipeline for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
