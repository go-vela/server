// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestPostgres_Client_Ping(t *testing.T) {
	// create the new fake sql database
	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new sql mock database: %v", err)
	}
	defer _sql.Close()

	// ensure the mock expects the ping
	_mock.ExpectPing()

	// setup the database client
	_database, err := NewTest(_sql)
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	// create the new closed fake SQL database
	_closedSQL, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new sql mock database: %v", err)
	}

	// setup the closed database client
	_closed, err := NewTest(_closedSQL)
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	// close the fake sql database to simulate failures to ping
	_closedSQL.Close()

	// setup tests
	tests := []struct {
		failure  bool
		database *client
	}{
		{
			failure:  false,
			database: _database,
		},
		{
			failure:  true,
			database: _closed,
		},
	}

	// run tests
	for _, test := range tests {
		err = test.database.Ping()

		if test.failure {
			if err == nil {
				t.Errorf("Ping should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Ping returned err: %v", err)
		}
	}
}
