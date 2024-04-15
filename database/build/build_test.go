// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/raw"
)

func TestBuild_New(t *testing.T) {
	// setup types
	logger := logrus.NewEntry(logrus.StandardLogger())

	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new SQL mock: %v", err)
	}
	defer _sql.Close()

	_mock.ExpectExec(CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(CreateCreatedIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(CreateRepoIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(CreateSourceIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(CreateStatusIndex).WillReturnResult(sqlmock.NewResult(1, 1))

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
				client: _postgres,
				config: &config{SkipCreation: false},
				ctx:    context.TODO(),
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
				client: _sqlite,
				config: &config{SkipCreation: false},
				ctx:    context.TODO(),
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

func TestBuild_Crop(t *testing.T) {
	// setup types
	title := randomString(1001)
	message := randomString(2001)
	err := randomString(1001)

	b := testBuild()
	b.Title = sql.NullString{String: title, Valid: true}
	b.Message = sql.NullString{String: message, Valid: true}
	b.Error = sql.NullString{String: err, Valid: true}

	want := testBuild()
	want.Title = sql.NullString{String: title[:1000], Valid: true}
	want.Message = sql.NullString{String: message[:2000], Valid: true}
	want.Error = sql.NullString{String: err[:1000], Valid: true}

	// run test
	got := b.Crop()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Crop is %v, want %v", got, want)
	}
}

func TestBuild_Nullify(t *testing.T) {
	// setup types
	var b *Build

	want := &Build{
		ID:            sql.NullInt64{Int64: 0, Valid: false},
		RepoID:        sql.NullInt64{Int64: 0, Valid: false},
		PipelineID:    sql.NullInt64{Int64: 0, Valid: false},
		Number:        sql.NullInt32{Int32: 0, Valid: false},
		Parent:        sql.NullInt32{Int32: 0, Valid: false},
		Event:         sql.NullString{String: "", Valid: false},
		EventAction:   sql.NullString{String: "", Valid: false},
		Status:        sql.NullString{String: "", Valid: false},
		Error:         sql.NullString{String: "", Valid: false},
		Enqueued:      sql.NullInt64{Int64: 0, Valid: false},
		Created:       sql.NullInt64{Int64: 0, Valid: false},
		Started:       sql.NullInt64{Int64: 0, Valid: false},
		Finished:      sql.NullInt64{Int64: 0, Valid: false},
		Deploy:        sql.NullString{String: "", Valid: false},
		DeployNumber:  sql.NullInt64{Int64: 0, Valid: false},
		DeployPayload: nil,
		Clone:         sql.NullString{String: "", Valid: false},
		Source:        sql.NullString{String: "", Valid: false},
		Title:         sql.NullString{String: "", Valid: false},
		Message:       sql.NullString{String: "", Valid: false},
		Commit:        sql.NullString{String: "", Valid: false},
		Sender:        sql.NullString{String: "", Valid: false},
		Author:        sql.NullString{String: "", Valid: false},
		Email:         sql.NullString{String: "", Valid: false},
		Link:          sql.NullString{String: "", Valid: false},
		Branch:        sql.NullString{String: "", Valid: false},
		Ref:           sql.NullString{String: "", Valid: false},
		BaseRef:       sql.NullString{String: "", Valid: false},
		HeadRef:       sql.NullString{String: "", Valid: false},
		Host:          sql.NullString{String: "", Valid: false},
		Runtime:       sql.NullString{String: "", Valid: false},
		Distribution:  sql.NullString{String: "", Valid: false},
	}

	// setup tests
	tests := []struct {
		build *Build
		want  *Build
	}{
		{
			build: testBuild(),
			want:  testBuild(),
		},
		{
			build: b,
			want:  nil,
		},
		{
			build: new(Build),
			want:  want,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.build.Nullify()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Nullify is %v, want %v", got, test.want)
		}
	}
}

func TestBuild_ToAPI(t *testing.T) {
	// setup types
	want := new(api.Build)

	want.SetID(1)
	want.SetRepo(testAPIRepo())
	want.SetPipelineID(1)
	want.SetNumber(1)
	want.SetParent(1)
	want.SetEvent("push")
	want.SetEventAction("")
	want.SetStatus("running")
	want.SetError("")
	want.SetEnqueued(1563474077)
	want.SetCreated(1563474076)
	want.SetStarted(1563474078)
	want.SetFinished(1563474079)
	want.SetDeploy("")
	want.SetDeployNumber(0)
	want.SetDeployPayload(nil)
	want.SetClone("https://github.com/github/octocat.git")
	want.SetSource("https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135163")
	want.SetTitle("push received from https://github.com/github/octocat")
	want.SetMessage("First commit...")
	want.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	want.SetSender("OctoKitty")
	want.SetAuthor("OctoKitty")
	want.SetEmail("OctoKitty@github.com")
	want.SetLink("https://example.company.com/github/octocat/1")
	want.SetBranch("main")
	want.SetRef("refs/heads/main")
	want.SetBaseRef("")
	want.SetHeadRef("")
	want.SetHost("example.company.com")
	want.SetRuntime("docker")
	want.SetDistribution("linux")
	want.SetDeployPayload(raw.StringSliceMap{"foo": "test1", "bar": "test2"})
	want.SetApprovedAt(1563474076)
	want.SetApprovedBy("OctoCat")

	// run test
	got := testBuild().ToAPI()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToAPI is %v, want %v", got, want)
	}
}

func TestBuild_Validate(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		build   *Build
	}{
		{
			failure: false,
			build:   testBuild(),
		},
		{ // no repo_id set for build
			failure: true,
			build: &Build{
				ID:     sql.NullInt64{Int64: 1, Valid: true},
				Number: sql.NullInt32{Int32: 1, Valid: true},
			},
		},
		{ // no number set for build
			failure: true,
			build: &Build{
				ID:     sql.NullInt64{Int64: 1, Valid: true},
				RepoID: sql.NullInt64{Int64: 1, Valid: true},
			},
		},
	}

	// run tests
	for _, test := range tests {
		err := test.build.Validate()

		if test.failure {
			if err == nil {
				t.Errorf("Validate should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Validate returned err: %v", err)
		}
	}
}

func TestBuild_FromAPI(t *testing.T) {
	// setup types
	b := new(api.Build)

	b.SetID(1)
	b.SetRepo(testAPIRepo())
	b.SetPipelineID(1)
	b.SetNumber(1)
	b.SetParent(1)
	b.SetEvent("push")
	b.SetEventAction("")
	b.SetStatus("running")
	b.SetError("")
	b.SetEnqueued(1563474077)
	b.SetCreated(1563474076)
	b.SetStarted(1563474078)
	b.SetFinished(1563474079)
	b.SetDeploy("")
	b.SetDeployNumber(0)
	b.SetDeployPayload(nil)
	b.SetClone("https://github.com/github/octocat.git")
	b.SetSource("https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135163")
	b.SetTitle("push received from https://github.com/github/octocat")
	b.SetMessage("First commit...")
	b.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	b.SetSender("OctoKitty")
	b.SetAuthor("OctoKitty")
	b.SetEmail("OctoKitty@github.com")
	b.SetLink("https://example.company.com/github/octocat/1")
	b.SetBranch("main")
	b.SetRef("refs/heads/main")
	b.SetBaseRef("")
	b.SetHeadRef("")
	b.SetHost("example.company.com")
	b.SetRuntime("docker")
	b.SetDistribution("linux")
	b.SetDeployPayload(raw.StringSliceMap{"foo": "test1", "bar": "test2"})
	b.SetApprovedAt(1563474076)
	b.SetApprovedBy("OctoCat")

	want := testBuild()

	// run test
	got := FromAPI(b)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("FromAPI is %v, want %v", got, want)
	}
}

func TestQueueBuild_ToAPI(t *testing.T) {
	// setup types
	want := new(api.QueueBuild)

	want.SetNumber(1)
	want.SetStatus("running")
	want.SetCreated(1563474076)
	want.SetFullName("github/octocat")

	// run test
	got := testQueueBuild().ToAPI()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToAPI is %v, want %v", got, want)
	}
}

func TestQueueBuild_FromAPI(t *testing.T) {
	// setup types
	b := new(api.QueueBuild)

	b.SetNumber(1)
	b.SetStatus("running")
	b.SetCreated(1563474076)
	b.SetFullName("github/octocat")

	want := testQueueBuild()

	// run test
	got := QueueBuildFromAPI(b)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("QueueBuildFromAPI is %v, want %v", got, want)
	}
}

// TEST RESOURCES

// This will be used with the github.com/DATA-DOG/go-sqlmock library to compare values
// that are otherwise not easily compared. These typically would be values generated
// before adding or updating them in the database.
//
// https://github.com/DATA-DOG/go-sqlmock#matching-arguments-like-timetime
type AnyArgument struct{}

// Match satisfies sqlmock.Argument interface.
func (a AnyArgument) Match(_ driver.Value) bool {
	return true
}

// NowTimestamp is used to test whether timestamps get updated correctly to the current time with lenience.
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

func randomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		//nolint:gosec // accepting weak RNG for test
		b[i] = letter[rand.Intn(len(letter))]
	}

	return string(b)
}

// testBuild is a test helper function to create a Build
// type with all fields set to a fake value.
func testBuild() *Build {
	return &Build{
		ID:            sql.NullInt64{Int64: 1, Valid: true},
		RepoID:        sql.NullInt64{Int64: 1, Valid: true},
		PipelineID:    sql.NullInt64{Int64: 1, Valid: true},
		Number:        sql.NullInt32{Int32: 1, Valid: true},
		Parent:        sql.NullInt32{Int32: 1, Valid: true},
		Event:         sql.NullString{String: "push", Valid: true},
		EventAction:   sql.NullString{String: "", Valid: false},
		Status:        sql.NullString{String: "running", Valid: true},
		Error:         sql.NullString{String: "", Valid: false},
		Enqueued:      sql.NullInt64{Int64: 1563474077, Valid: true},
		Created:       sql.NullInt64{Int64: 1563474076, Valid: true},
		Started:       sql.NullInt64{Int64: 1563474078, Valid: true},
		Finished:      sql.NullInt64{Int64: 1563474079, Valid: true},
		Deploy:        sql.NullString{String: "", Valid: false},
		DeployNumber:  sql.NullInt64{Int64: 0, Valid: false},
		DeployPayload: raw.StringSliceMap{"foo": "test1", "bar": "test2"},
		Clone:         sql.NullString{String: "https://github.com/github/octocat.git", Valid: true},
		Source:        sql.NullString{String: "https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135163", Valid: true},
		Title:         sql.NullString{String: "push received from https://github.com/github/octocat", Valid: true},
		Message:       sql.NullString{String: "First commit...", Valid: true},
		Commit:        sql.NullString{String: "48afb5bdc41ad69bf22588491333f7cf71135163", Valid: true},
		Sender:        sql.NullString{String: "OctoKitty", Valid: true},
		Author:        sql.NullString{String: "OctoKitty", Valid: true},
		Email:         sql.NullString{String: "OctoKitty@github.com", Valid: true},
		Link:          sql.NullString{String: "https://example.company.com/github/octocat/1", Valid: true},
		Branch:        sql.NullString{String: "main", Valid: true},
		Ref:           sql.NullString{String: "refs/heads/main", Valid: true},
		BaseRef:       sql.NullString{String: "", Valid: false},
		HeadRef:       sql.NullString{String: "", Valid: false},
		Host:          sql.NullString{String: "example.company.com", Valid: true},
		Runtime:       sql.NullString{String: "docker", Valid: true},
		Distribution:  sql.NullString{String: "linux", Valid: true},
		ApprovedAt:    sql.NullInt64{Int64: 1563474076, Valid: true},
		ApprovedBy:    sql.NullString{String: "OctoCat", Valid: true},
	}
}

// testQueueBuild is a test helper function to create a QueueBuild
// type with all fields set to a fake value.
func testQueueBuild() *QueueBuild {
	return &QueueBuild{
		Number:   sql.NullInt32{Int32: 1, Valid: true},
		Status:   sql.NullString{String: "running", Valid: true},
		Created:  sql.NullInt64{Int64: 1563474076, Valid: true},
		FullName: sql.NullString{String: "github/octocat", Valid: true},
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
	_mock.ExpectExec(CreateCreatedIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(CreateRepoIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(CreateSourceIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(CreateStatusIndex).WillReturnResult(sqlmock.NewResult(1, 1))

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
		t.Errorf("unable to create new postgres build engine: %v", err)
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
		t.Errorf("unable to create new sqlite build engine: %v", err)
	}

	return _engine
}

// testAPIBuild is a test helper function to create a API
// Build type with all fields set to their zero values.
func testAPIBuild() *api.Build {
	return &api.Build{
		ID:           new(int64),
		Repo:         testAPIRepo(),
		PipelineID:   new(int64),
		Number:       new(int),
		Parent:       new(int),
		Event:        new(string),
		EventAction:  new(string),
		Status:       new(string),
		Error:        new(string),
		Enqueued:     new(int64),
		Created:      new(int64),
		Started:      new(int64),
		Finished:     new(int64),
		Deploy:       new(string),
		DeployNumber: new(int64),
		Clone:        new(string),
		Source:       new(string),
		Title:        new(string),
		Message:      new(string),
		Commit:       new(string),
		Sender:       new(string),
		Author:       new(string),
		Email:        new(string),
		Link:         new(string),
		Branch:       new(string),
		Ref:          new(string),
		BaseRef:      new(string),
		HeadRef:      new(string),
		Host:         new(string),
		Runtime:      new(string),
		Distribution: new(string),
		ApprovedAt:   new(int64),
		ApprovedBy:   new(string),
	}
}

// testDeployment is a test helper function to create a library
// Repo type with all fields set to their zero values.
func testDeployment() *library.Deployment {
	builds := []*library.Build{}
	return &library.Deployment{
		ID:          new(int64),
		RepoID:      new(int64),
		Number:      new(int64),
		URL:         new(string),
		Commit:      new(string),
		Ref:         new(string),
		Task:        new(string),
		Target:      new(string),
		Description: new(string),
		Payload:     new(raw.StringSliceMap),
		CreatedAt:   new(int64),
		CreatedBy:   new(string),
		Builds:      builds,
	}
}

// testAPIRepo is a test helper function to create an API
// Repo type with all fields set to their zero values.
func testAPIRepo() *api.Repo {
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

func testAPIUser() *api.User {
	return &api.User{
		ID:           new(int64),
		Name:         new(string),
		RefreshToken: new(string),
		Token:        new(string),
		Favorites:    new([]string),
		Active:       new(bool),
		Admin:        new(bool),
	}
}
