// SPDX-License-Identifier: Apache-2.0

package executable

import (
	"database/sql/driver"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

func TestExecutable_New(t *testing.T) {
	// setup types
	logger := logrus.NewEntry(logrus.StandardLogger())

	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new SQL mock: %v", err)
	}
	defer _sql.Close()

	_mock.ExpectExec(CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))

	_config := &gorm.Config{SkipDefaultTransaction: true}

	_postgres, err := gorm.Open(postgres.New(postgres.Config{Conn: _sql}), _config)
	if err != nil {
		t.Errorf("unable to create new postgres database: %v", err)
	}

	_sqlite, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), _config)
	if err != nil {
		t.Errorf("unable to create new sqlite database: %v", err)
	}

	defer func() { _sql, _ := _sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure      bool
		name         string
		client       *gorm.DB
		level        int
		key          string
		logger       *logrus.Entry
		skipCreation bool
		want         *Engine
	}{
		{
			failure:      false,
			name:         "postgres",
			client:       _postgres,
			level:        1,
			key:          "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
			logger:       logger,
			skipCreation: false,
			want: &Engine{
				client: _postgres,
				config: &config{
					CompressionLevel: 1,
					EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
					SkipCreation:     false,
					Driver:           "postgres",
				},
				logger: logger,
			},
		},
		{
			failure:      false,
			name:         "sqlite3",
			client:       _sqlite,
			level:        1,
			key:          "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
			logger:       logger,
			skipCreation: false,
			want: &Engine{
				client: _sqlite,
				config: &config{
					CompressionLevel: 1,
					EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
					SkipCreation:     false,
					Driver:           "sqlite3",
				},
				logger: logger,
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := New(
				WithClient(test.client),
				WithCompressionLevel(test.level),
				WithEncryptionKey(test.key),
				WithLogger(test.logger),
				WithSkipCreation(test.skipCreation),
				WithDriver(test.name),
			)

			if test.failure {
				if err == nil {
					t.Errorf("New for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("New for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("New for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}

// testPostgres is a helper function to create a Postgres engine for testing.
func testPostgres(t *testing.T) (*Engine, sqlmock.Sqlmock) {
	// create the new mock sql database
	//
	// https://pkg.go.dev/github.com/DATA-DOG/go-sqlmock#New
	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new SQL mock: %v", err)
	}

	_mock.ExpectExec(CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))

	// create the new mock Postgres database client
	//
	// https://pkg.go.dev/gorm.io/gorm#Open
	_postgres, err := gorm.Open(
		postgres.New(postgres.Config{Conn: _sql}),
		&gorm.Config{SkipDefaultTransaction: true},
	)
	if err != nil {
		t.Errorf("unable to create new postgres database: %v", err)
	}

	_engine, err := New(
		WithClient(_postgres),
		WithCompressionLevel(0),
		WithEncryptionKey("A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW"),
		WithLogger(logrus.NewEntry(logrus.StandardLogger())),
		WithSkipCreation(false),
		WithDriver(constants.DriverPostgres),
	)
	if err != nil {
		t.Errorf("unable to create new postgres build_itnerary engine: %v", err)
	}

	return _engine, _mock
}

// testSqlite is a helper function to create a Sqlite engine for testing.
func testSqlite(t *testing.T) *Engine {
	_sqlite, err := gorm.Open(
		sqlite.Open("file::memory:?cache=shared"),
		&gorm.Config{SkipDefaultTransaction: true},
	)
	if err != nil {
		t.Errorf("unable to create new sqlite database: %v", err)
	}

	_engine, err := New(
		WithClient(_sqlite),
		WithCompressionLevel(0),
		WithEncryptionKey("A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW"),
		WithLogger(logrus.NewEntry(logrus.StandardLogger())),
		WithSkipCreation(false),
		WithDriver(constants.DriverSqlite),
	)
	if err != nil {
		t.Errorf("unable to create new sqlite build_itnerary engine: %v", err)
	}

	return _engine
}

// testBuildExecutable is a test helper function to create a API
// BuildExecutable type with all fields set to their zero values.
func testBuildExecutable() *api.BuildExecutable {
	return &api.BuildExecutable{
		ID:      new(int64),
		BuildID: new(int64),
		Data:    new([]byte),
	}
}

// This will be used with the github.com/DATA-DOG/go-sqlmock library to compare values
// that are otherwise not easily compared. These typically would be values generated
// before adding or updating them in the database.
//
// https://github.com/DATA-DOG/go-sqlmock#matching-arguments-like-timetime
type AnyArgument struct{}

// Match satisfies sqlmock.Argument interface.
func (a AnyArgument) Match(v driver.Value) bool {
	return v != nil
}
