// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

func TestDatabase_Log_Compress(t *testing.T) {
	// setup tests
	tests := []struct {
		name    string
		failure bool
		level   int
		log     *Log
		want    []byte
	}{
		{
			name:    "compression level -1",
			failure: false,
			level:   constants.CompressionNegOne,
			log:     &Log{Data: []byte("foo")},
			want:    []byte{120, 156, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
		},
		{
			name:    "compression level 0",
			failure: false,
			level:   constants.CompressionZero,
			log:     &Log{Data: []byte("foo")},
			want:    []byte{120, 1, 0, 3, 0, 252, 255, 102, 111, 111, 1, 0, 0, 255, 255, 2, 130, 1, 69},
		},
		{
			name:    "compression level 1",
			failure: false,
			level:   constants.CompressionOne,
			log:     &Log{Data: []byte("foo")},
			want:    []byte{120, 1, 0, 3, 0, 252, 255, 102, 111, 111, 1, 0, 0, 255, 255, 2, 130, 1, 69},
		},
		{
			name:    "compression level 2",
			failure: false,
			level:   constants.CompressionTwo,
			log:     &Log{Data: []byte("foo")},
			want:    []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
		},
		{
			name:    "compression level 3",
			failure: false,
			level:   constants.CompressionThree,
			log:     &Log{Data: []byte("foo")},
			want:    []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
		},
		{
			name:    "compression level 4",
			failure: false,
			level:   constants.CompressionFour,
			log:     &Log{Data: []byte("foo")},
			want:    []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
		},
		{
			name:    "compression level 5",
			failure: false,
			level:   constants.CompressionFive,
			log:     &Log{Data: []byte("foo")},
			want:    []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
		},
		{
			name:    "compression level 6",
			failure: false,
			level:   constants.CompressionSix,
			log:     &Log{Data: []byte("foo")},
			want:    []byte{120, 156, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
		},
		{
			name:    "compression level 7",
			failure: false,
			level:   constants.CompressionSeven,
			log:     &Log{Data: []byte("foo")},
			want:    []byte{120, 218, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
		},
		{
			name:    "compression level 8",
			failure: false,
			level:   constants.CompressionEight,
			log:     &Log{Data: []byte("foo")},
			want:    []byte{120, 218, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
		},
		{
			name:    "compression level 9",
			failure: false,
			level:   constants.CompressionNine,
			log:     &Log{Data: []byte("foo")},
			want:    []byte{120, 218, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.log.Compress(test.level)

			if test.failure {
				if err == nil {
					t.Errorf("Compress for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("Compress for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(test.log.Data, test.want) {
				t.Errorf("Compress for %s is %v, want %v", test.name, string(test.log.Data), string(test.want))
			}
		})
	}
}

func TestDatabase_Log_Decompress(t *testing.T) {
	// setup tests
	tests := []struct {
		name    string
		failure bool
		log     *Log
		want    []byte
	}{
		{
			name:    "compression level -1",
			failure: false,
			log:     &Log{Data: []byte{120, 156, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69}},
			want:    []byte("foo"),
		},
		{
			name:    "compression level 0",
			failure: false,
			log:     &Log{Data: []byte{120, 1, 0, 3, 0, 252, 255, 102, 111, 111, 1, 0, 0, 255, 255, 2, 130, 1, 69}},
			want:    []byte("foo"),
		},
		{
			name:    "compression level 1",
			failure: false,
			log:     &Log{Data: []byte{120, 1, 0, 3, 0, 252, 255, 102, 111, 111, 1, 0, 0, 255, 255, 2, 130, 1, 69}},
			want:    []byte("foo"),
		},
		{
			name:    "compression level 2",
			failure: false,
			log:     &Log{Data: []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69}},
			want:    []byte("foo"),
		},
		{
			name:    "compression level 3",
			failure: false,
			log:     &Log{Data: []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69}},
			want:    []byte("foo"),
		},
		{
			name:    "compression level 4",
			failure: false,
			log:     &Log{Data: []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69}},
			want:    []byte("foo"),
		},
		{
			name:    "compression level 5",
			failure: false,
			log:     &Log{Data: []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69}},
			want:    []byte("foo"),
		},
		{
			name:    "compression level 6",
			failure: false,
			log:     &Log{Data: []byte{120, 156, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69}},
			want:    []byte("foo"),
		},
		{
			name:    "compression level 7",
			failure: false,
			log:     &Log{Data: []byte{120, 218, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69}},
			want:    []byte("foo"),
		},
		{
			name:    "compression level 8",
			failure: false,
			log:     &Log{Data: []byte{120, 218, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69}},
			want:    []byte("foo"),
		},
		{
			name:    "compression level 9",
			failure: false,
			log:     &Log{Data: []byte{120, 218, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69}},
			want:    []byte("foo"),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.log.Decompress()

			if test.failure {
				if err == nil {
					t.Errorf("Decompress for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("Decompress for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(test.log.Data, test.want) {
				t.Errorf("Decompress for %s is %v, want %v", test.name, string(test.log.Data), string(test.want))
			}
		})
	}
}

func TestDatabase_Log_Nullify(t *testing.T) {
	// setup types
	var l *Log

	want := &Log{
		ID:        sql.NullInt64{Int64: 0, Valid: false},
		BuildID:   sql.NullInt64{Int64: 0, Valid: false},
		RepoID:    sql.NullInt64{Int64: 0, Valid: false},
		ServiceID: sql.NullInt64{Int64: 0, Valid: false},
		StepID:    sql.NullInt64{Int64: 0, Valid: false},
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
	}
}
