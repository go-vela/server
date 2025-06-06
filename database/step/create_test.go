// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestStep_Engine_CreateStep(t *testing.T) {
	// setup types
	_step := testutils.APIStep()
	_step.SetID(1)
	_step.SetRepoID(1)
	_step.SetBuildID(1)
	_step.SetNumber(1)
	_step.SetName("foo")
	_step.SetImage("bar")
	_step.SetReportAs("test")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "steps"
("build_id","repo_id","number","name","image","stage","status","error","exit_code","created","started","finished","host","runtime","distribution","report_as","id")
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17) RETURNING "id"`).
		WithArgs(1, 1, 1, "foo", "bar", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, "test", 1).
		WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.CreateStep(context.TODO(), _step)

			if test.failure {
				if err == nil {
					t.Errorf("CreateStep for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateStep for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _step) {
				t.Errorf("CreateStep for %s returned %s, want %s", test.name, got, _step)
			}
		})
	}
}
