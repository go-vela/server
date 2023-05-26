// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package service

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestService_Engine_CleanService(t *testing.T) {
	// setup types
	_serviceOne := testService()
	_serviceOne.SetID(1)
	_serviceOne.SetRepoID(1)
	_serviceOne.SetBuildID(1)
	_serviceOne.SetNumber(1)
	_serviceOne.SetName("foo")
	_serviceOne.SetImage("bar")
	_serviceOne.SetCreated(1)
	_serviceOne.SetStatus("running")

	_serviceTwo := testService()
	_serviceTwo.SetID(2)
	_serviceTwo.SetRepoID(1)
	_serviceTwo.SetBuildID(1)
	_serviceTwo.SetNumber(2)
	_serviceTwo.SetName("foo")
	_serviceTwo.SetImage("bar")
	_serviceTwo.SetCreated(1)
	_serviceTwo.SetStatus("pending")

	_serviceThree := testService()
	_serviceThree.SetID(3)
	_serviceThree.SetRepoID(1)
	_serviceThree.SetBuildID(1)
	_serviceThree.SetNumber(3)
	_serviceThree.SetName("foo")
	_serviceThree.SetImage("bar")
	_serviceThree.SetCreated(1)
	_serviceThree.SetStatus("success")

	_serviceFour := testService()
	_serviceFour.SetID(4)
	_serviceFour.SetRepoID(1)
	_serviceFour.SetBuildID(1)
	_serviceFour.SetNumber(4)
	_serviceFour.SetName("foo")
	_serviceFour.SetImage("bar")
	_serviceFour.SetCreated(5)
	_serviceFour.SetStatus("pending")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the name query
	_mock.ExpectExec(`UPDATE "services" SET "status"=$1 WHERE created < $2 AND (status = 'running' OR status = 'pending')`).
		WithArgs("error", 3).
		WillReturnResult(sqlmock.NewResult(1, 2))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateService(_serviceOne)
	if err != nil {
		t.Errorf("unable to create test service for sqlite: %v", err)
	}

	err = _sqlite.CreateService(_serviceTwo)
	if err != nil {
		t.Errorf("unable to create test service for sqlite: %v", err)
	}

	err = _sqlite.CreateService(_serviceThree)
	if err != nil {
		t.Errorf("unable to create test service for sqlite: %v", err)
	}

	err = _sqlite.CreateService(_serviceFour)
	if err != nil {
		t.Errorf("unable to create test service for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
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
			got, err := test.database.CleanServices(3)

			if test.failure {
				if err == nil {
					t.Errorf("CleanServices for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CleanServices for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("CleanServices for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
