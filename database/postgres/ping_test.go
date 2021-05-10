// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"testing"
)

func TestPostgres_Client_Ping(t *testing.T) {
	// setup types
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() {
		_sql, _ := _database.Postgres.DB()
		_sql.Close()
	}()

	_bad, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	// close the bad database to simulate failures to ping
	_sql, _ := _bad.Postgres.DB()
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
			database: _bad,
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
