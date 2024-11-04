// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
)

func TestDatabase_Step_Nullify(t *testing.T) {
	// setup types
	var s *Step

	want := &Step{
		ID:           sql.NullInt64{Int64: 0, Valid: false},
		BuildID:      sql.NullInt64{Int64: 0, Valid: false},
		RepoID:       sql.NullInt64{Int64: 0, Valid: false},
		Number:       sql.NullInt64{Int64: 0, Valid: false},
		Name:         sql.NullString{String: "", Valid: false},
		Image:        sql.NullString{String: "", Valid: false},
		Stage:        sql.NullString{String: "", Valid: false},
		Status:       sql.NullString{String: "", Valid: false},
		Error:        sql.NullString{String: "", Valid: false},
		ExitCode:     sql.NullInt64{Int64: 0, Valid: false},
		Created:      sql.NullInt64{Int64: 0, Valid: false},
		Started:      sql.NullInt64{Int64: 0, Valid: false},
		Finished:     sql.NullInt64{Int64: 0, Valid: false},
		Host:         sql.NullString{String: "", Valid: false},
		Runtime:      sql.NullString{String: "", Valid: false},
		Distribution: sql.NullString{String: "", Valid: false},
		ReportAs:     sql.NullString{String: "", Valid: false},
	}

	// setup tests
	tests := []struct {
		step *Step
		want *Step
	}{
		{
			step: testStep(),
			want: testStep(),
		},
		{
			step: s,
			want: nil,
		},
		{
			step: new(Step),
			want: want,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.step.Nullify()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Nullify is %v, want %v", got, test.want)
		}
	}
}

func TestDatabase_Step_ToAPI(t *testing.T) {
	// setup types
	want := new(api.Step)

	want.SetID(1)
	want.SetBuildID(1)
	want.SetRepoID(1)
	want.SetNumber(1)
	want.SetName("clone")
	want.SetImage("target/vela-git:v0.3.0")
	want.SetStage("")
	want.SetStatus("running")
	want.SetError("")
	want.SetExitCode(0)
	want.SetCreated(1563474076)
	want.SetStarted(1563474078)
	want.SetFinished(1563474079)
	want.SetHost("example.company.com")
	want.SetRuntime("docker")
	want.SetDistribution("linux")
	want.SetReportAs("test")

	// run test
	got := testStep().ToAPI()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToAPI is %v, want %v", got, want)
	}
}

func TestDatabase_Step_Validate(t *testing.T) {
	// setup types
	tests := []struct {
		failure bool
		step    *Step
	}{
		{
			failure: false,
			step:    testStep(),
		},
		{ // no build_id set for step
			failure: true,
			step: &Step{
				ID:     sql.NullInt64{Int64: 1, Valid: true},
				RepoID: sql.NullInt64{Int64: 1, Valid: true},
				Number: sql.NullInt64{Int64: 1, Valid: true},
				Name:   sql.NullString{String: "clone", Valid: true},
				Image:  sql.NullString{String: "target/vela-git:v0.3.0", Valid: true},
			},
		},
		{ // no repo_id set for step
			failure: true,
			step: &Step{
				ID:      sql.NullInt64{Int64: 1, Valid: true},
				BuildID: sql.NullInt64{Int64: 1, Valid: true},
				Number:  sql.NullInt64{Int64: 1, Valid: true},
				Name:    sql.NullString{String: "clone", Valid: true},
				Image:   sql.NullString{String: "target/vela-git:v0.3.0", Valid: true},
			},
		},
		{ // no name set for step
			failure: true,
			step: &Step{
				ID:      sql.NullInt64{Int64: 1, Valid: true},
				BuildID: sql.NullInt64{Int64: 1, Valid: true},
				RepoID:  sql.NullInt64{Int64: 1, Valid: true},
				Number:  sql.NullInt64{Int64: 1, Valid: true},
				Image:   sql.NullString{String: "target/vela-git:v0.3.0", Valid: true},
			},
		},
		{ // no number set for step
			failure: true,
			step: &Step{
				ID:      sql.NullInt64{Int64: 1, Valid: true},
				BuildID: sql.NullInt64{Int64: 1, Valid: true},
				RepoID:  sql.NullInt64{Int64: 1, Valid: true},
				Name:    sql.NullString{String: "clone", Valid: true},
				Image:   sql.NullString{String: "target/vela-git:v0.3.0", Valid: true},
			},
		},
		{ // no image set for step
			failure: true,
			step: &Step{
				ID:      sql.NullInt64{Int64: 1, Valid: true},
				BuildID: sql.NullInt64{Int64: 1, Valid: true},
				RepoID:  sql.NullInt64{Int64: 1, Valid: true},
				Number:  sql.NullInt64{Int64: 1, Valid: true},
				Name:    sql.NullString{String: "clone", Valid: true},
			},
		},
	}

	// run tests
	for _, test := range tests {
		err := test.step.Validate()

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

func TestDatabase_StepFromAPI(t *testing.T) {
	// setup types
	s := new(api.Step)

	s.SetID(1)
	s.SetBuildID(1)
	s.SetRepoID(1)
	s.SetNumber(1)
	s.SetName("clone")
	s.SetImage("target/vela-git:v0.3.0")
	s.SetStage("")
	s.SetStatus("running")
	s.SetError("")
	s.SetExitCode(0)
	s.SetCreated(1563474076)
	s.SetStarted(1563474078)
	s.SetFinished(1563474079)
	s.SetHost("example.company.com")
	s.SetRuntime("docker")
	s.SetDistribution("linux")
	s.SetReportAs("test")

	want := testStep()

	// run test
	got := StepFromAPI(s)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("StepFromAPI is %v, want %v", got, want)
	}
}

// testStep is a test helper function to create a Step
// type with all fields set to a fake value.
func testStep() *Step {
	return &Step{
		ID:           sql.NullInt64{Int64: 1, Valid: true},
		BuildID:      sql.NullInt64{Int64: 1, Valid: true},
		RepoID:       sql.NullInt64{Int64: 1, Valid: true},
		Number:       sql.NullInt64{Int64: 1, Valid: true},
		Name:         sql.NullString{String: "clone", Valid: true},
		Image:        sql.NullString{String: "target/vela-git:v0.3.0", Valid: true},
		Stage:        sql.NullString{String: "", Valid: false},
		Status:       sql.NullString{String: "running", Valid: true},
		Error:        sql.NullString{String: "", Valid: false},
		ExitCode:     sql.NullInt64{Int64: 0, Valid: false},
		Created:      sql.NullInt64{Int64: 1563474076, Valid: true},
		Started:      sql.NullInt64{Int64: 1563474078, Valid: true},
		Finished:     sql.NullInt64{Int64: 1563474079, Valid: true},
		Host:         sql.NullString{String: "example.company.com", Valid: true},
		Runtime:      sql.NullString{String: "docker", Valid: true},
		Distribution: sql.NullString{String: "linux", Valid: true},
		ReportAs:     sql.NullString{String: "test", Valid: true},
	}
}
