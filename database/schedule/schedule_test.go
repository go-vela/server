// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	api "github.com/go-vela/server/api/types"
)

func TestSchedule_New(t *testing.T) {
	// setup types
	logger := logrus.NewEntry(logrus.StandardLogger())

	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new SQL mock: %v", err)
	}
	defer _sql.Close()

	_mock.ExpectExec(CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(CreateRepoIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))

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
		key          string
		logger       *logrus.Entry
		skipCreation bool
		want         *engine
	}{
		{
			failure:      false,
			name:         "postgres",
			client:       _postgres,
			logger:       logger,
			skipCreation: false,
			want: &engine{
				ctx:    context.TODO(),
				client: _postgres,
				config: &config{SkipCreation: false},
				logger: logger,
			},
		},
		{
			failure:      false,
			name:         "sqlite3",
			client:       _sqlite,
			logger:       logger,
			skipCreation: false,
			want: &engine{
				ctx:    context.TODO(),
				client: _sqlite,
				config: &config{SkipCreation: false},
				logger: logger,
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := New(
				WithContext(context.TODO()),
				WithClient(test.client),
				WithLogger(test.logger),
				WithSkipCreation(test.skipCreation),
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
func testPostgres(t *testing.T) (*engine, sqlmock.Sqlmock) {
	// create the new mock sql database
	//
	// https://pkg.go.dev/github.com/DATA-DOG/go-sqlmock#New
	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new SQL mock: %v", err)
	}

	_mock.ExpectExec(CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(CreateRepoIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))

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
		WithContext(context.TODO()),
		WithClient(_postgres),
		WithLogger(logrus.NewEntry(logrus.StandardLogger())),
		WithSkipCreation(false),
	)
	if err != nil {
		t.Errorf("unable to create new postgres schedule engine: %v", err)
	}

	return _engine, _mock
}

// testSqlite is a helper function to create a Sqlite engine for testing.
func testSqlite(t *testing.T) *engine {
	_sqlite, err := gorm.Open(
		sqlite.Open("file::memory:?cache=shared"),
		&gorm.Config{SkipDefaultTransaction: true},
	)
	if err != nil {
		t.Errorf("unable to create new sqlite database: %v", err)
	}

	_engine, err := New(
		WithContext(context.TODO()),
		WithClient(_sqlite),
		WithLogger(logrus.NewEntry(logrus.StandardLogger())),
		WithSkipCreation(false),
	)
	if err != nil {
		t.Errorf("unable to create new sqlite schedule engine: %v", err)
	}

	return _engine
}

// testSchedule is a test helper function to create an API Schedule type with all fields set to their zero values.
func testAPISchedule() *api.Schedule {
	return &api.Schedule{
		ID:          new(int64),
		RepoID:      new(int64),
		Active:      new(bool),
		Name:        new(string),
		Entry:       new(string),
		CreatedAt:   new(int64),
		CreatedBy:   new(string),
		UpdatedAt:   new(int64),
		UpdatedBy:   new(string),
		ScheduledAt: new(int64),
		Branch:      new(string),
		Error:       new(string),
	}
}

// This will be used with the github.com/DATA-DOG/go-sqlmock library to compare values
// that are otherwise not easily compared. These typically would be values generated
// before adding or updating them in the database.
//
// https://github.com/DATA-DOG/go-sqlmock#matching-arguments-like-timetime
type NowTimestamp struct{}

// Match satisfies sqlmock.Argument interface.
func (t NowTimestamp) Match(v driver.Value) bool {
	ts, ok := v.(int64)
	if !ok {
		return false
	}
	now := time.Now().Unix()

	return now-ts < 10
}

func TestSchedule_FromAPI(t *testing.T) {
	s := new(api.Schedule)
	s.SetID(1)
	s.SetRepoID(1)
	s.SetActive(true)
	s.SetName("nightly")
	s.SetEntry("0 0 * * *")
	s.SetCreatedAt(time.Now().UTC().Unix())
	s.SetCreatedBy("user1")
	s.SetUpdatedAt(time.Now().Add(time.Hour * 1).UTC().Unix())
	s.SetUpdatedBy("user2")
	s.SetScheduledAt(time.Now().Add(time.Hour * 2).UTC().Unix())
	s.SetBranch("main")
	s.SetError("unable to trigger build for schedule nightly: unknown character")

	want := testSchedule()

	got := FromAPI(s)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ScheduleFromAPI is %v, want %v", got, want)
	}
}

func TestSchedule_Nullify(t *testing.T) {
	tests := []struct {
		name     string
		schedule *Schedule
		want     *Schedule
	}{
		{
			name:     "schedule with fields",
			schedule: testSchedule(),
			want:     testSchedule(),
		},
		{
			name:     "schedule with empty fields",
			schedule: new(Schedule),
			want: &Schedule{
				ID:          sql.NullInt64{Int64: 0, Valid: false},
				RepoID:      sql.NullInt64{Int64: 0, Valid: false},
				Active:      sql.NullBool{Bool: false, Valid: false},
				Name:        sql.NullString{String: "", Valid: false},
				Entry:       sql.NullString{String: "", Valid: false},
				CreatedAt:   sql.NullInt64{Int64: 0, Valid: false},
				CreatedBy:   sql.NullString{String: "", Valid: false},
				UpdatedAt:   sql.NullInt64{Int64: 0, Valid: false},
				UpdatedBy:   sql.NullString{String: "", Valid: false},
				ScheduledAt: sql.NullInt64{Int64: 0, Valid: false},
				Branch:      sql.NullString{String: "", Valid: false},
				Error:       sql.NullString{String: "", Valid: false},
			},
		},
		{
			name:     "empty schedule",
			schedule: nil,
			want:     nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.schedule.Nullify()
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Nullify is %v, want %v", got, test.want)
			}
		})
	}
}

func TestSchedule_ToAPI(t *testing.T) {
	want := new(api.Schedule)
	want.SetID(1)
	want.SetRepoID(1)
	want.SetActive(true)
	want.SetName("nightly")
	want.SetEntry("0 0 * * *")
	want.SetCreatedAt(time.Now().UTC().Unix())
	want.SetCreatedBy("user1")
	want.SetUpdatedAt(time.Now().Add(time.Hour * 1).UTC().Unix())
	want.SetUpdatedBy("user2")
	want.SetScheduledAt(time.Now().Add(time.Hour * 2).UTC().Unix())
	want.SetBranch("main")
	want.SetError("unable to trigger build for schedule nightly: unknown character")

	got := testSchedule().ToAPI()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToLibrary is %v, want %v", got, want)
	}
}

func TestSchedule_Validate(t *testing.T) {
	tests := []struct {
		name     string
		failure  bool
		schedule *Schedule
	}{
		{
			name:     "schedule with valid fields",
			failure:  false,
			schedule: testSchedule(),
		},
		{
			name:    "schedule with invalid entry",
			failure: true,
			schedule: &Schedule{
				ID:     sql.NullInt64{Int64: 1, Valid: true},
				RepoID: sql.NullInt64{Int64: 1, Valid: true},
				Name:   sql.NullString{String: "invalid", Valid: false},
				Entry:  sql.NullString{String: "!@#$%^&*()", Valid: false},
			},
		},
		{
			name:    "schedule with missing entry",
			failure: true,
			schedule: &Schedule{
				ID:     sql.NullInt64{Int64: 1, Valid: true},
				RepoID: sql.NullInt64{Int64: 1, Valid: true},
				Name:   sql.NullString{String: "nightly", Valid: false},
			},
		},
		{
			name:    "schedule with missing name",
			failure: true,
			schedule: &Schedule{
				ID:     sql.NullInt64{Int64: 1, Valid: true},
				RepoID: sql.NullInt64{Int64: 1, Valid: true},
				Entry:  sql.NullString{String: "0 0 * * *", Valid: false},
			},
		},
		{
			name:    "schedule with missing repo_id",
			failure: true,
			schedule: &Schedule{
				ID:    sql.NullInt64{Int64: 1, Valid: true},
				Name:  sql.NullString{String: "nightly", Valid: false},
				Entry: sql.NullString{String: "0 0 * * *", Valid: false},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.schedule.Validate()
			if test.failure {
				if err == nil {
					t.Errorf("Validate should have returned err")
				}

				return
			}

			if err != nil {
				t.Errorf("Validate returned err: %v", err)
			}
		})
	}
}

// testSchedule is a test helper function to create a Schedule type with all fields set to a fake value.
func testSchedule() *Schedule {
	return &Schedule{
		ID:          sql.NullInt64{Int64: 1, Valid: true},
		RepoID:      sql.NullInt64{Int64: 1, Valid: true},
		Active:      sql.NullBool{Bool: true, Valid: true},
		Name:        sql.NullString{String: "nightly", Valid: true},
		Entry:       sql.NullString{String: "0 0 * * *", Valid: true},
		CreatedAt:   sql.NullInt64{Int64: time.Now().UTC().Unix(), Valid: true},
		CreatedBy:   sql.NullString{String: "user1", Valid: true},
		UpdatedAt:   sql.NullInt64{Int64: time.Now().Add(time.Hour * 1).UTC().Unix(), Valid: true},
		UpdatedBy:   sql.NullString{String: "user2", Valid: true},
		ScheduledAt: sql.NullInt64{Int64: time.Now().Add(time.Hour * 2).UTC().Unix(), Valid: true},
		Branch:      sql.NullString{String: "main", Valid: true},
		Error:       sql.NullString{String: "unable to trigger build for schedule nightly: unknown character", Valid: true},
	}
}

// testRepo is a test helper function to create a library Repo type with all fields set to their zero values.
func testRepo() *api.Repo {
	return &api.Repo{
		ID:           new(int64),
		BuildLimit:   new(int64),
		Timeout:      new(int64),
		Counter:      new(int),
		PipelineType: new(string),
		Hash:         new(string),
		Org:          new(string),
		Name:         new(string),
		FullName:     new(string),
		Link:         new(string),
		Clone:        new(string),
		Branch:       new(string),
		Visibility:   new(string),
		PreviousName: new(string),
		Private:      new(bool),
		Trusted:      new(bool),
		Active:       new(bool),
	}
}
