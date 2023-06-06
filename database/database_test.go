// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDatabase_New(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		name    string
		config  *Config
	}{
		{
			name:    "failure with postgres",
			failure: true,
			config: &Config{
				Driver:           "postgres",
				Address:          "postgres://foo:bar@localhost:5432/vela",
				CompressionLevel: 3,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
				SkipCreation:     false,
			},
		},
		{
			name:    "success with sqlite3",
			failure: false,
			config: &Config{
				Driver:           "sqlite3",
				Address:          "file::memory:?cache=shared",
				CompressionLevel: 3,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
				SkipCreation:     false,
			},
		},
		{
			name:    "failure with invalid config",
			failure: true,
			config: &Config{
				Driver:           "postgres",
				Address:          "",
				CompressionLevel: 3,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
				SkipCreation:     false,
			},
		},
		{
			name:    "failure with invalid driver",
			failure: true,
			config: &Config{
				Driver:           "mysql",
				Address:          "foo:bar@tcp(localhost:3306)/vela?charset=utf8mb4&parseTime=True&loc=Local",
				CompressionLevel: 3,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
				SkipCreation:     false,
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := New(test.config)

			if test.failure {
				if err == nil {
					t.Errorf("New for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("New for %s returned err: %v", test.name, err)
			}
		})
	}
}

// testPostgres is a helper function to create a Postgres engine for testing.
func testPostgres(t *testing.T) (*engine, sqlmock.Sqlmock) {
	// create the engine with test configuration
	_engine := &engine{
		Config: &Config{
			CompressionLevel: 3,
			ConnectionLife:   30 * time.Minute,
			ConnectionIdle:   2,
			ConnectionOpen:   0,
			Driver:           "postgres",
			EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
			SkipCreation:     false,
		},
		Logger: logrus.NewEntry(logrus.StandardLogger()),
	}

	// create the new mock sql database
	_sql, _mock, err := sqlmock.New(
		sqlmock.MonitorPingsOption(true),
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	if err != nil {
		t.Errorf("unable to create new SQL mock: %v", err)
	}
	// ensure the mock expects the ping
	_mock.ExpectPing()

	// create the new mock Postgres database client
	_engine.Database, err = gorm.Open(
		postgres.New(postgres.Config{Conn: _sql}),
		&gorm.Config{SkipDefaultTransaction: true},
	)
	if err != nil {
		t.Errorf("unable to create new test postgres database: %v", err)
	}

	return _engine, _mock
}

// testSqlite is a helper function to create a Sqlite engine for testing.
func testSqlite(t *testing.T) *engine {
	var err error

	// create the engine with test configuration
	_engine := &engine{
		Config: &Config{
			Address:          "file::memory:?cache=shared",
			CompressionLevel: 3,
			ConnectionLife:   30 * time.Minute,
			ConnectionIdle:   2,
			ConnectionOpen:   0,
			Driver:           "sqlite3",
			EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
			SkipCreation:     false,
		},
		Logger: logrus.NewEntry(logrus.StandardLogger()),
	}

	// create the new mock Sqlite database client
	_engine.Database, err = gorm.Open(
		sqlite.Open(_engine.Config.Address),
		&gorm.Config{SkipDefaultTransaction: true},
	)
	if err != nil {
		t.Errorf("unable to create new test sqlite database: %v", err)
	}

	return _engine
}
