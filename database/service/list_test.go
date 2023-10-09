// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestService_Engine_ListServices(t *testing.T) {
	// setup types
	_serviceOne := testService()
	_serviceOne.SetID(1)
	_serviceOne.SetRepoID(1)
	_serviceOne.SetBuildID(1)
	_serviceOne.SetNumber(1)
	_serviceOne.SetName("foo")
	_serviceOne.SetImage("bar")

	_serviceTwo := testService()
	_serviceTwo.SetID(2)
	_serviceTwo.SetRepoID(1)
	_serviceTwo.SetBuildID(2)
	_serviceTwo.SetNumber(1)
	_serviceTwo.SetName("bar")
	_serviceTwo.SetImage("foo")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "services"`).WillReturnRows(_rows)

	// create expected result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "name", "image", "stage", "status", "error", "exit_code", "created", "started", "finished", "host", "runtime", "distribution"}).
		AddRow(1, 1, 1, 1, "foo", "bar", "", "", "", 0, 0, 0, 0, "", "", "").
		AddRow(2, 1, 2, 1, "bar", "foo", "", "", "", 0, 0, 0, 0, "", "", "")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "services"`).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateService(context.TODO(), _serviceOne)
	if err != nil {
		t.Errorf("unable to create test service for sqlite: %v", err)
	}

	_, err = _sqlite.CreateService(context.TODO(), _serviceTwo)
	if err != nil {
		t.Errorf("unable to create test service for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     []*library.Service
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*library.Service{_serviceOne, _serviceTwo},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*library.Service{_serviceOne, _serviceTwo},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListServices(context.TODO())

			if test.failure {
				if err == nil {
					t.Errorf("ListServices for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListServices for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListServices for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
