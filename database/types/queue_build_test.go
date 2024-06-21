// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
)

func TestTypes_QueueBuild_ToAPI(t *testing.T) {
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

func TestTypes_QueueBuild_FromAPI(t *testing.T) {
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
