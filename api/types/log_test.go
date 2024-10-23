// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTypes_Log_AppendData(t *testing.T) {
	// setup types
	data := []byte("bar")

	want := testLog()
	want.SetData([]byte("foobar"))

	// setup tests
	tests := []struct {
		log  *Log
		want *Log
	}{
		{
			log:  testLog(),
			want: want,
		},
		{
			log:  new(Log),
			want: &Log{Data: &data},
		},
	}

	// run tests
	for _, test := range tests {
		test.log.AppendData(data)

		if !reflect.DeepEqual(test.log, test.want) {
			t.Errorf("AppendData is %v, want %v", test.log, test.want)
		}
	}
}

func TestTypes_Log_MaskData(t *testing.T) {
	// set up test secrets
	sVals := []string{"gh_abc123def456", "((%.YY245***pP.><@@}}", "quick-bear-fox-squid", "SUPERSECRETVALUE"}

	tests := []struct {
		log     []byte
		want    []byte
		secrets []string
	}{
		{ // no secrets in log
			log: []byte(
				"$ echo hello\nhello\n",
			),
			want: []byte(
				"$ echo hello\nhello\n",
			),
			secrets: sVals,
		},
		{ // one secret in log
			log: []byte(
				"((%.YY245***pP.><@@}}",
			),
			want: []byte(
				"***",
			),
			secrets: sVals,
		},
		{ // multiple secrets in log
			log: []byte(
				"$ echo $SECRET1\n((%.YY245***pP.><@@}}\n$ echo $SECRET2\nquick-bear-fox-squid\n",
			),
			want: []byte(
				"$ echo $SECRET1\n***\n$ echo $SECRET2\n***\n",
			),
			secrets: sVals,
		},
		{ // secret with leading =
			log: []byte(
				"SOME_SECRET=((%.YY245***pP.><@@}}",
			),
			want: []byte(
				"SOME_SECRET=***",
			),
			secrets: sVals,
		},
		{ // secret baked in URL query params
			log: []byte(
				"www.example.com?username=quick-bear-fox-squid&password=SUPERSECRETVALUE",
			),
			want: []byte(
				"www.example.com?username=***&password=***",
			),
			secrets: sVals,
		},
		{ // secret in verbose brackets
			log: []byte(
				"[token: gh_abc123def456]",
			),
			want: []byte(
				"[token: ***]",
			),
			secrets: sVals,
		},
		{ // double secret
			log: []byte(
				"echo ${GITHUB_TOKEN}${SUPER_SECRET}\ngh_abc123def456SUPERSECRETVALUE\n",
			),
			want: []byte(
				"echo ${GITHUB_TOKEN}${SUPER_SECRET}\n******\n",
			),
			secrets: sVals,
		},
		{ // empty secrets slice
			log: []byte(
				"echo hello\nhello\n",
			),
			want: []byte(
				"echo hello\nhello\n",
			),
			secrets: []string{},
		},
	}
	// run tests
	l := testLog()
	for _, test := range tests {
		l.SetData(test.log)
		l.MaskData(test.secrets)

		got := l.GetData()
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("MaskData is %v, want %v", string(got), string(test.want))
		}
	}
}

func TestTypes_Log_Getters(t *testing.T) {
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
			log:  new(Log),
			want: new(Log),
		},
	}

	// run tests
	for _, test := range tests {
		if test.log.GetID() != test.want.GetID() {
			t.Errorf("GetID is %v, want %v", test.log.GetID(), test.want.GetID())
		}

		if test.log.GetServiceID() != test.want.GetServiceID() {
			t.Errorf("GetServiceID is %v, want %v", test.log.GetServiceID(), test.want.GetServiceID())
		}

		if test.log.GetStepID() != test.want.GetStepID() {
			t.Errorf("GetStepID is %v, want %v", test.log.GetStepID(), test.want.GetStepID())
		}

		if test.log.GetBuildID() != test.want.GetBuildID() {
			t.Errorf("GetBuildID is %v, want %v", test.log.GetBuildID(), test.want.GetBuildID())
		}

		if test.log.GetRepoID() != test.want.GetRepoID() {
			t.Errorf("GetRepoID is %v, want %v", test.log.GetRepoID(), test.want.GetRepoID())
		}

		if !reflect.DeepEqual(test.log.GetData(), test.want.GetData()) {
			t.Errorf("GetData is %v, want %v", test.log.GetData(), test.want.GetData())
		}
	}
}

func TestTypes_Log_Setters(t *testing.T) {
	// setup types
	var l *Log

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
			want: new(Log),
		},
	}

	// run tests
	for _, test := range tests {
		test.log.SetID(test.want.GetID())
		test.log.SetServiceID(test.want.GetServiceID())
		test.log.SetStepID(test.want.GetStepID())
		test.log.SetBuildID(test.want.GetBuildID())
		test.log.SetRepoID(test.want.GetRepoID())
		test.log.SetData(test.want.GetData())

		if test.log.GetID() != test.want.GetID() {
			t.Errorf("SetID is %v, want %v", test.log.GetID(), test.want.GetID())
		}

		if test.log.GetServiceID() != test.want.GetServiceID() {
			t.Errorf("SetServiceID is %v, want %v", test.log.GetServiceID(), test.want.GetServiceID())
		}

		if test.log.GetStepID() != test.want.GetStepID() {
			t.Errorf("SetStepID is %v, want %v", test.log.GetStepID(), test.want.GetStepID())
		}

		if test.log.GetBuildID() != test.want.GetBuildID() {
			t.Errorf("SetBuildID is %v, want %v", test.log.GetBuildID(), test.want.GetBuildID())
		}

		if test.log.GetRepoID() != test.want.GetRepoID() {
			t.Errorf("SetRepoID is %v, want %v", test.log.GetRepoID(), test.want.GetRepoID())
		}

		if !reflect.DeepEqual(test.log.GetData(), test.want.GetData()) {
			t.Errorf("SetData is %v, want %v", test.log.GetData(), test.want.GetData())
		}
	}
}

func TestTypes_Log_String(t *testing.T) {
	// setup types
	l := testLog()

	want := fmt.Sprintf(`{
  BuildID: %d,
  Data: %s,
  ID: %d,
  RepoID: %d,
  ServiceID: %d,
  StepID: %d,
}`,
		l.GetBuildID(),
		l.GetData(),
		l.GetID(),
		l.GetRepoID(),
		l.GetServiceID(),
		l.GetStepID(),
	)

	// run test
	got := l.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

// testLog is a test helper function to create a Log
// type with all fields set to a fake value.
func testLog() *Log {
	l := new(Log)

	l.SetID(1)
	l.SetServiceID(1)
	l.SetStepID(1)
	l.SetBuildID(1)
	l.SetRepoID(1)
	l.SetData([]byte("foo"))

	return l
}
