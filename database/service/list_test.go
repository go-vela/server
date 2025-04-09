// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestService_Engine_ListServices(t *testing.T) {
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
	_serviceTwo.SetBuildID(2)
	_serviceTwo.SetNumber(1)
	_serviceTwo.SetName("bar")
	_serviceTwo.SetImage("foo")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.ServiceFromAPI(_serviceOne), *types.ServiceFromAPI(_serviceTwo)})

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
		database *Engine
		want     []*api.Service
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*api.Service{_serviceOne, _serviceTwo},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.Service{_serviceOne, _serviceTwo},
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
