// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
)

func TestDatabase_Log_Nullify(t *testing.T) {
	// setup types
	var l *Log

	want := &Log{
		ID:        sql.NullInt64{Int64: 0, Valid: false},
		BuildID:   sql.NullInt64{Int64: 0, Valid: false},
		RepoID:    sql.NullInt64{Int64: 0, Valid: false},
		ServiceID: sql.NullInt64{Int64: 0, Valid: false},
		StepID:    sql.NullInt64{Int64: 0, Valid: false},
		CreatedAt: sql.NullInt64{Int64: 0, Valid: false},
	}

	// setup tests
	tests := []struct {
		log  *Log
		want *Log
	}{
		{
			log:  testLog(),
			want: testLog(),
		},
		{
			log:  l,
			want: nil,
		},
		{
			log:  new(Log),
			want: want,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.log.Nullify()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Nullify is %v, want %v", got, test.want)
		}
	}
}

func TestDatabase_Log_ToAPI(t *testing.T) {
	// setup types
	want := new(api.Log)

	want.SetID(1)
	want.SetServiceID(1)
	want.SetStepID(1)
	want.SetBuildID(1)
	want.SetRepoID(1)
	want.SetData([]byte("foo"))
	want.SetCreatedAt(tsCreate)

	// run test
	got := testLog().ToAPI()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToAPI is %v, want %v", got, want)
	}
}

func TestDatabase_Log_Validate(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		log     *Log
	}{
		{
			failure: false,
			log:     testLog(),
		},
		{ // no service_id or step_id set for log
			failure: true,
			log: &Log{
				ID:      sql.NullInt64{Int64: 1, Valid: true},
				BuildID: sql.NullInt64{Int64: 1, Valid: true},
				RepoID:  sql.NullInt64{Int64: 1, Valid: true},
			},
		},
		{ // no build_id set for log
			failure: true,
			log: &Log{
				ID:        sql.NullInt64{Int64: 1, Valid: true},
				RepoID:    sql.NullInt64{Int64: 1, Valid: true},
				ServiceID: sql.NullInt64{Int64: 1, Valid: true},
				StepID:    sql.NullInt64{Int64: 1, Valid: true},
			},
		},
		{ // no repo_id set for log
			failure: true,
			log: &Log{
				ID:        sql.NullInt64{Int64: 1, Valid: true},
				BuildID:   sql.NullInt64{Int64: 1, Valid: true},
				ServiceID: sql.NullInt64{Int64: 1, Valid: true},
				StepID:    sql.NullInt64{Int64: 1, Valid: true},
			},
		},
	}

	// run tests
	for _, test := range tests {
		err := test.log.Validate()

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

func TestDatabase_LogFromAPI(t *testing.T) {
	// setup types
	l := new(api.Log)

	l.SetID(1)
	l.SetServiceID(1)
	l.SetStepID(1)
	l.SetBuildID(1)
	l.SetRepoID(1)
	l.SetData([]byte("foo"))
	l.SetCreatedAt(tsCreate)

	want := testLog()

	// run test
	got := LogFromAPI(l)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("LogFromAPI is %v, want %v", got, want)
	}
}

// testLog is a test helper function to create a Log
// type with all fields set to a fake value.
func testLog() *Log {
	return &Log{
		ID:        sql.NullInt64{Int64: 1, Valid: true},
		BuildID:   sql.NullInt64{Int64: 1, Valid: true},
		RepoID:    sql.NullInt64{Int64: 1, Valid: true},
		ServiceID: sql.NullInt64{Int64: 1, Valid: true},
		StepID:    sql.NullInt64{Int64: 1, Valid: true},
		Data:      []byte("foo"),
		CreatedAt: sql.NullInt64{Int64: tsCreate, Valid: true},
	}
}
