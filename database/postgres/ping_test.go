// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"testing"
)

func TestPostgres_Client_Ping(t *testing.T) {
	// setup types
	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// ensure the mock expects the ping
	_mock.ExpectPing()

	// setup the closed test database client
	_closed, _, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	// capture the closed test sql database
	_sql, _ := _closed.Postgres.DB()
	// close the test sql database to simulate failures to ping
	_sql.Close()

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
