// SPDX-License-Identifier: Apache-2.0

package executable

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestExecutable_Engine_CleanExecutables(t *testing.T) {
	// setup types
	_buildOne := testutils.APIBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepo(testutils.APIRepo())
	_buildOne.SetNumber(1)
	_buildOne.SetStatus("pending")
	_buildOne.SetCreated(1)
	_buildOne.SetDeployPayload(nil)

	_buildTwo := testutils.APIBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepo(testutils.APIRepo())
	_buildTwo.SetNumber(2)
	_buildTwo.SetStatus("error")
	_buildTwo.SetCreated(1)
	_buildTwo.SetDeployPayload(nil)

	_bExecutableOne := testBuildExecutable()
	_bExecutableOne.SetID(1)
	_bExecutableOne.SetBuildID(1)
	_bExecutableOne.SetData([]byte("foo"))

	_bExecutableTwo := testBuildExecutable()
	_bExecutableTwo.SetID(2)
	_bExecutableTwo.SetBuildID(2)
	_bExecutableTwo.SetData([]byte("bar"))

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	_mock.ExpectExec("DELETE FROM build_executables USING builds WHERE builds.id = build_executables.build_id AND builds.status = 'error';").
		WillReturnResult(sqlmock.NewResult(1, 1))

	_mock.ExpectQuery(`DELETE FROM "build_executables" WHERE build_id = $1 RETURNING *`).WithArgs(2).WillReturnError(fmt.Errorf("not found"))

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateBuildExecutable(context.TODO(), _bExecutableOne)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	err = _sqlite.CreateBuildExecutable(context.TODO(), _bExecutableTwo)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	err = _sqlite.client.AutoMigrate(&types.Build{})
	if err != nil {
		t.Errorf("unable to create repo table for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableBuild).Create(types.BuildFromAPI(_buildOne)).Error
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableBuild).Create(types.BuildFromAPI(_buildTwo)).Error
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

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
			got, err := test.database.CleanBuildExecutables(context.TODO())

			if test.failure {
				if err == nil {
					t.Errorf("CleanExecutables for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CleanExecutables for %s returned err: %v", test.name, err)
			}

			if got != 1 {
				t.Errorf("CleanExecutables for %s should have affected 1 row, affected %d", test.name, got)
			}

			_, err = test.database.PopBuildExecutable(context.TODO(), 2)
			if err == nil {
				t.Errorf("CleanExecutables for %s should have returned an error", test.name)
			}
		})
	}
}
