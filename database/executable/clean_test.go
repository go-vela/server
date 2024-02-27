// SPDX-License-Identifier: Apache-2.0

package executable

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

func TestExecutable_Engine_CleanExecutables(t *testing.T) {
	// setup types
	_buildOne := testBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepoID(1)
	_buildOne.SetNumber(1)
	_buildOne.SetStatus("pending")
	_buildOne.SetCreated(1)
	_buildOne.SetDeployPayload(nil)

	_buildTwo := testBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepoID(1)
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

	err = _sqlite.client.AutoMigrate(&database.Build{})
	if err != nil {
		t.Errorf("unable to create repo table for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableBuild).Create(database.BuildFromLibrary(_buildOne)).Error
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableBuild).Create(database.BuildFromLibrary(_buildTwo)).Error
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
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

// testBuild is a test helper function to create a library
// Build type with all fields set to their zero values.
func testBuild() *library.Build {
	return &library.Build{
		ID:           new(int64),
		RepoID:       new(int64),
		PipelineID:   new(int64),
		Number:       new(int),
		Parent:       new(int),
		Event:        new(string),
		EventAction:  new(string),
		Status:       new(string),
		Error:        new(string),
		Enqueued:     new(int64),
		Created:      new(int64),
		Started:      new(int64),
		Finished:     new(int64),
		Deploy:       new(string),
		Clone:        new(string),
		Source:       new(string),
		Title:        new(string),
		Message:      new(string),
		Commit:       new(string),
		Sender:       new(string),
		Author:       new(string),
		Email:        new(string),
		Link:         new(string),
		Branch:       new(string),
		Ref:          new(string),
		BaseRef:      new(string),
		HeadRef:      new(string),
		Host:         new(string),
		Runtime:      new(string),
		Distribution: new(string),
	}
}
