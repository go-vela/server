// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestStep_Engine_ListStepsForBuild(t *testing.T) {
	// setup types
	_build := testutils.APIBuild()
	_build.SetID(1)
	_build.SetRepo(testutils.APIRepo())
	_build.SetNumber(1)

	_stepOne := testutils.APIStep()
	_stepOne.SetID(1)
	_stepOne.SetRepoID(1)
	_stepOne.SetBuildID(1)
	_stepOne.SetNumber(1)
	_stepOne.SetName("foo")
	_stepOne.SetImage("bar")

	_stepTwo := testutils.APIStep()
	_stepTwo.SetID(2)
	_stepTwo.SetRepoID(1)
	_stepTwo.SetBuildID(1)
	_stepTwo.SetNumber(2)
	_stepTwo.SetName("foo")
	_stepTwo.SetImage("bar")
	_stepTwo.SetReportAs("test")

	_postgres, _mock := testPostgres(t)

	ctx := context.TODO()
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.StepFromAPI(_stepTwo), *types.StepFromAPI(_stepOne)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "steps" WHERE build_id = $1 ORDER BY id DESC LIMIT $2`).WithArgs(1, 10).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateStep(ctx, _stepOne)
	if err != nil {
		t.Errorf("unable to create test step for sqlite: %v", err)
	}

	_, err = _sqlite.CreateStep(ctx, _stepTwo)
	if err != nil {
		t.Errorf("unable to create test step for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     []*api.Step
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*api.Step{_stepTwo, _stepOne},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.Step{_stepTwo, _stepOne},
		},
	}

	filters := map[string]interface{}{}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListStepsForBuild(ctx, _build, filters, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListStepsForBuild for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListStepsForBuild for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListStepsForBuild for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
