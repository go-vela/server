// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
)

func TestDatabase_BuildExecutable_Decrypt(t *testing.T) {
	// setup types
	key := "C639A572E14D5075C526FDDD43E4ECF6"
	encrypted := testBuildExecutable()

	err := encrypted.Encrypt(key)
	if err != nil {
		t.Errorf("unable to encrypt repo: %v", err)
	}

	// setup tests
	tests := []struct {
		failure    bool
		key        string
		executable BuildExecutable
	}{
		{
			failure:    false,
			key:        key,
			executable: *encrypted,
		},
		{
			failure:    true,
			key:        "",
			executable: *encrypted,
		},
		{
			failure:    true,
			key:        key,
			executable: *testBuildExecutable(),
		},
	}

	// run tests
	for _, test := range tests {
		err := test.executable.Decrypt(test.key)

		if test.failure {
			if err == nil {
				t.Errorf("Decrypt should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Decrypt returned err: %v", err)
		}
	}
}

func TestDatabase_BuildExecutable_Encrypt(t *testing.T) {
	// setup types
	key := "C639A572E14D5075C526FDDD43E4ECF6"

	// setup tests
	tests := []struct {
		failure    bool
		key        string
		executable *BuildExecutable
	}{
		{
			failure:    false,
			key:        key,
			executable: testBuildExecutable(),
		},
		{
			failure:    true,
			key:        "",
			executable: testBuildExecutable(),
		},
	}

	// run tests
	for _, test := range tests {
		err := test.executable.Encrypt(test.key)

		if test.failure {
			if err == nil {
				t.Errorf("Encrypt should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Encrypt returned err: %v", err)
		}
	}
}

func TestDatabase_BuildExecutable_Nullify(t *testing.T) {
	// setup types
	var p *BuildExecutable

	want := &BuildExecutable{
		ID:      sql.NullInt64{Int64: 0, Valid: false},
		BuildID: sql.NullInt64{Int64: 0, Valid: false},
	}

	// setup tests
	tests := []struct {
		buildExecutable *BuildExecutable
		want            *BuildExecutable
	}{
		{
			buildExecutable: testBuildExecutable(),
			want:            testBuildExecutable(),
		},
		{
			buildExecutable: p,
			want:            nil,
		},
		{
			buildExecutable: new(BuildExecutable),
			want:            want,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.buildExecutable.Nullify()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Nullify is %v, want %v", got, test.want)
		}
	}
}

func TestDatabase_BuildExecutable_ToAPI(t *testing.T) {
	// setup types
	want := new(api.BuildExecutable)

	want.SetID(1)
	want.SetBuildID(1)
	want.SetData(testBuildExecutableData())

	// run test
	got := testBuildExecutable().ToAPI()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToAPI is %v, want %v", got, want)
	}
}

func TestDatabase_BuildExecutable_Validate(t *testing.T) {
	// setup tests
	tests := []struct {
		failure         bool
		buildExecutable *BuildExecutable
	}{
		{
			failure:         false,
			buildExecutable: testBuildExecutable(),
		},
		{ // no build_id set for buildExecutable
			failure: true,
			buildExecutable: &BuildExecutable{
				ID: sql.NullInt64{Int64: 1, Valid: true},
			},
		},
	}

	// run tests
	for _, test := range tests {
		err := test.buildExecutable.Validate()

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

func TestDatabase_BuildExecutableFromAPI(t *testing.T) {
	// setup types
	c := new(api.BuildExecutable)

	c.SetID(1)
	c.SetBuildID(1)
	c.SetData(testBuildExecutableData())

	want := testBuildExecutable()

	// run test
	got := BuildExecutableFromAPI(c)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("BuildExecutableFromAPI is %v, want %v", got, want)
	}
}

// testBuildExecutable is a test helper function to create a BuildExecutable
// type with all fields set to a fake value.
func testBuildExecutable() *BuildExecutable {
	return &BuildExecutable{
		ID:      sql.NullInt64{Int64: 1, Valid: true},
		BuildID: sql.NullInt64{Int64: 1, Valid: true},
		Data:    testBuildExecutableData(),
	}
}

// testBuildExecutableData is a test helper function to create the
// content for the Data field for the BuildExecutable type.
func testBuildExecutableData() []byte {
	return []byte(`
{ 
    "id": "step_name",
    "version": "1",
    "metadata":{
        "clone":true,
        "environment":["steps","services","secrets"]},
    "worker":{},
    "steps":[
        {
            "id":"step_github_octocat_1_init",
            "directory":"/vela/src/github.com/github/octocat",
            "environment": {"BUILD_AUTHOR":"Octocat"}
        }
    ]
}
`)
}
