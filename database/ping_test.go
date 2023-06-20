// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"testing"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func TestDatabase_Engine_Ping(t *testing.T) {
	_postgres, _mock := testPostgres(t)
	defer _postgres.Close()
	// ensure the mock expects the ping
	_mock.ExpectPing()

	// create a test database without mocking the call
	_unmocked, _ := testPostgres(t)

	_sqlite := testSqlite(t)
	defer _sqlite.Close()

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
	}{
		{
			name:     "success with postgres",
			failure:  false,
			database: _postgres,
		},
		{
			name:     "success with sqlite",
			failure:  false,
			database: _sqlite,
		},
		{
			name:     "failure without mocked call",
			failure:  true,
			database: _unmocked,
		},
		{
			name:    "failure with invalid gorm database",
			failure: true,
			database: &engine{
				config: &config{
					Driver: "invalid",
				},
				client: &gorm.DB{
					Config: &gorm.Config{
						ConnPool: nil,
					},
				},
				logger: logrus.NewEntry(logrus.StandardLogger()),
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.database.Ping()

			if test.failure {
				if err == nil {
					t.Errorf("Ping for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("Ping for %s returned err: %v", test.name, err)
			}
		})
	}
}
