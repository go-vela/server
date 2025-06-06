// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestService_Engine_ListServiceStatusCount(t *testing.T) {
	// setup types
	_serviceOne := testutils.APIService()
	_serviceOne.SetID(1)
	_serviceOne.SetRepoID(1)
	_serviceOne.SetBuildID(1)
	_serviceOne.SetNumber(1)
	_serviceOne.SetName("foo")
	_serviceOne.SetImage("bar")

	_serviceTwo := testutils.APIService()
	_serviceTwo.SetID(2)
	_serviceTwo.SetRepoID(1)
	_serviceTwo.SetBuildID(1)
	_serviceTwo.SetNumber(2)
	_serviceTwo.SetName("foo")
	_serviceTwo.SetImage("bar")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"status", "count"}).
		AddRow("pending", 0).
		AddRow("failure", 0).
		AddRow("killed", 0).
		AddRow("running", 0).
		AddRow("success", 0)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT "status", count(status) as count FROM "services" GROUP BY "status"`).WillReturnRows(_rows)

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
		database *Engine
		want     map[string]float64
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want: map[string]float64{
				"pending": 0,
				"failure": 0,
				"killed":  0,
				"running": 0,
				"success": 0,
			},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want: map[string]float64{
				"pending": 0,
				"failure": 0,
				"killed":  0,
				"running": 0,
				"success": 0,
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListServiceStatusCount(context.TODO())

			if test.failure {
				if err == nil {
					t.Errorf("ListServiceStatusCount for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListServiceStatusCount for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListServiceStatusCount for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
