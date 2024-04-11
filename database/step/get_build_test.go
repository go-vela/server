// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/types/library"
)

func TestStep_Engine_GetStepForBuild(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)

	_step := testStep()
	_step.SetID(1)
	_step.SetRepoID(1)
	_step.SetBuildID(1)
	_step.SetNumber(1)
	_step.SetName("foo")
	_step.SetImage("bar")

	_postgres, _mock := testPostgres(t)

	ctx := context.TODO()

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "name", "image", "stage", "status", "error", "exit_code", "created", "started", "finished", "host", "runtime", "distribution", "report_as"}).
		AddRow(1, 1, 1, 1, "foo", "bar", "", "", "", 0, 0, 0, 0, "", "", "", "")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "steps" WHERE build_id = $1 AND number = $2 LIMIT $3`).WithArgs(1, 1, 1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()
	_, err := _sqlite.CreateStep(ctx, _step)
	if err != nil {
		t.Errorf("unable to create test step for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     *library.Step
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _step,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _step,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetStepForBuild(ctx, _build, 1)

			if test.failure {
				if err == nil {
					t.Errorf("GetStepForBuild for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetStepForBuild for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetStepForBuild for %s is a mismatch (-want +got):\n%s", test.name, diff)
			}
		})
	}
}
