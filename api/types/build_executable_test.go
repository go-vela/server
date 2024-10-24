// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTypes_BuildExecutable_Getters(t *testing.T) {
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
			buildExecutable: new(BuildExecutable),
			want:            new(BuildExecutable),
		},
	}

	// run tests
	for _, test := range tests {
		if test.buildExecutable.GetID() != test.want.GetID() {
			t.Errorf("GetID is %v, want %v", test.buildExecutable.GetID(), test.want.GetID())
		}

		if test.buildExecutable.GetBuildID() != test.want.GetBuildID() {
			t.Errorf("GetBuildID is %v, want %v", test.buildExecutable.GetBuildID(), test.want.GetBuildID())
		}

		if !reflect.DeepEqual(test.buildExecutable.GetData(), test.want.GetData()) {
			t.Errorf("GetData is %v, want %v", test.buildExecutable.GetData(), test.want.GetData())
		}
	}
}

func TestTypes_BuildExecutable_Setters(t *testing.T) {
	// setup types
	var bExecutable *BuildExecutable

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
			buildExecutable: bExecutable,
			want:            new(BuildExecutable),
		},
	}

	// run tests
	for _, test := range tests {
		test.buildExecutable.SetID(test.want.GetID())
		test.buildExecutable.SetBuildID(test.want.GetBuildID())
		test.buildExecutable.SetData(test.want.GetData())

		if test.buildExecutable.GetID() != test.want.GetID() {
			t.Errorf("SetID is %v, want %v", test.buildExecutable.GetID(), test.want.GetID())
		}

		if test.buildExecutable.GetBuildID() != test.want.GetBuildID() {
			t.Errorf("SetRepoID is %v, want %v", test.buildExecutable.GetBuildID(), test.want.GetBuildID())
		}

		if !reflect.DeepEqual(test.buildExecutable.GetData(), test.want.GetData()) {
			t.Errorf("SetData is %v, want %v", test.buildExecutable.GetData(), test.want.GetData())
		}
	}
}

func TestTypes_BuildExecutable_String(t *testing.T) {
	// setup types
	bExecutable := testBuildExecutable()

	want := fmt.Sprintf(`{
  ID: %d,
  BuildID: %d,
  Data: %s,
}`,
		bExecutable.GetID(),
		bExecutable.GetBuildID(),
		bExecutable.GetData(),
	)

	// run test
	got := bExecutable.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

// testBuildExecutable is a test helper function to create a Pipeline
// type with all fields set to a fake value.
func testBuildExecutable() *BuildExecutable {
	p := new(BuildExecutable)

	p.SetID(1)
	p.SetBuildID(1)
	p.SetData(testBuildExecutableData())

	return p
}

// testBuildExecutableData is a test helper function to create the
// content for the Data field for the Pipeline type.
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
