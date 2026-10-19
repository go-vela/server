// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

func TestDatabase_Pipeline_Nullify(t *testing.T) {
	// setup types
	var p *Pipeline

	want := &Pipeline{
		ID:       sql.NullInt64{Int64: 0, Valid: false},
		RepoID:   sql.NullInt64{Int64: 0, Valid: false},
		Commit:   sql.NullString{String: "", Valid: false},
		Flavor:   sql.NullString{String: "", Valid: false},
		Platform: sql.NullString{String: "", Valid: false},
		Ref:      sql.NullString{String: "", Valid: false},
		Type:     sql.NullString{String: "", Valid: false},
		Version:  sql.NullString{String: "", Valid: false},
	}

	// setup tests
	tests := []struct {
		pipeline *Pipeline
		want     *Pipeline
	}{
		{
			pipeline: testPipeline(),
			want:     testPipeline(),
		},
		{
			pipeline: p,
			want:     nil,
		},
		{
			pipeline: new(Pipeline),
			want:     want,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.pipeline.Nullify()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Nullify is %v, want %v", got, test.want)
		}
	}
}

func TestDatabase_Pipeline_ToAPI(t *testing.T) {
	// setup types
	want := new(api.Pipeline)

	want.SetID(1)
	want.SetRepo(testRepo().ToAPI())
	want.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	want.SetFlavor("large")
	want.SetPlatform("docker")
	want.SetRef("refs/heads/main")
	want.SetType(constants.PipelineTypeYAML)
	want.SetVersion("1")
	want.SetExternalSecrets(false)
	want.SetInternalSecrets(false)
	want.SetServices(true)
	want.SetStages(false)
	want.SetSteps(true)
	want.SetTemplates(false)
	want.SetArtifact(false)
	want.SetWarnings([]string{"42:this is a warning"})
	want.SetData(testPipelineData())

	// run test
	got := testPipeline().ToAPI()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToAPI is %v, want %v", got, want)
	}
}

func TestDatabase_Pipeline_Validate(t *testing.T) {
	// setup tests
	tests := []struct {
		failure  bool
		pipeline *Pipeline
	}{
		{
			failure:  false,
			pipeline: testPipeline(),
		},
		{ // no commit set for pipeline
			failure: true,
			pipeline: &Pipeline{
				ID:      sql.NullInt64{Int64: 1, Valid: true},
				RepoID:  sql.NullInt64{Int64: 1, Valid: true},
				Ref:     sql.NullString{String: "refs/heads/main", Valid: true},
				Type:    sql.NullString{String: constants.PipelineTypeYAML, Valid: true},
				Version: sql.NullString{String: "1", Valid: true},
			},
		},
		{ // no ref set for pipeline
			failure: true,
			pipeline: &Pipeline{
				ID:      sql.NullInt64{Int64: 1, Valid: true},
				RepoID:  sql.NullInt64{Int64: 1, Valid: true},
				Commit:  sql.NullString{String: "48afb5bdc41ad69bf22588491333f7cf71135163", Valid: true},
				Type:    sql.NullString{String: constants.PipelineTypeYAML, Valid: true},
				Version: sql.NullString{String: "1", Valid: true},
			},
		},
		{ // no repo_id set for pipeline
			failure: true,
			pipeline: &Pipeline{
				ID:      sql.NullInt64{Int64: 1, Valid: true},
				Commit:  sql.NullString{String: "48afb5bdc41ad69bf22588491333f7cf71135163", Valid: true},
				Ref:     sql.NullString{String: "refs/heads/main", Valid: true},
				Type:    sql.NullString{String: constants.PipelineTypeYAML, Valid: true},
				Version: sql.NullString{String: "1", Valid: true},
			},
		},
		{ // no type set for pipeline
			failure: true,
			pipeline: &Pipeline{
				ID:      sql.NullInt64{Int64: 1, Valid: true},
				RepoID:  sql.NullInt64{Int64: 1, Valid: true},
				Commit:  sql.NullString{String: "48afb5bdc41ad69bf22588491333f7cf71135163", Valid: true},
				Ref:     sql.NullString{String: "refs/heads/main", Valid: true},
				Version: sql.NullString{String: "1", Valid: true},
			},
		},
		{ // no version set for pipeline
			failure: true,
			pipeline: &Pipeline{
				ID:     sql.NullInt64{Int64: 1, Valid: true},
				RepoID: sql.NullInt64{Int64: 1, Valid: true},
				Commit: sql.NullString{String: "48afb5bdc41ad69bf22588491333f7cf71135163", Valid: true},
				Ref:    sql.NullString{String: "refs/heads/main", Valid: true},

				Type: sql.NullString{String: constants.PipelineTypeYAML, Valid: true},
			},
		},
	}

	// run tests
	for _, test := range tests {
		err := test.pipeline.Validate()

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

func TestDatabase_PipelineFromAPI(t *testing.T) {
	// setup types
	want := &Pipeline{
		ID:              sql.NullInt64{Int64: 1, Valid: true},
		RepoID:          sql.NullInt64{Int64: 1, Valid: true},
		Commit:          sql.NullString{String: "48afb5bdc41ad69bf22588491333f7cf71135163", Valid: true},
		Flavor:          sql.NullString{String: "large", Valid: true},
		Platform:        sql.NullString{String: "docker", Valid: true},
		Ref:             sql.NullString{String: "refs/heads/main", Valid: true},
		Type:            sql.NullString{String: constants.PipelineTypeYAML, Valid: true},
		Version:         sql.NullString{String: "1", Valid: true},
		ExternalSecrets: sql.NullBool{Bool: false, Valid: true},
		InternalSecrets: sql.NullBool{Bool: false, Valid: true},
		Services:        sql.NullBool{Bool: true, Valid: true},
		Stages:          sql.NullBool{Bool: false, Valid: true},
		Steps:           sql.NullBool{Bool: true, Valid: true},
		Templates:       sql.NullBool{Bool: false, Valid: true},
		Artifacts:       sql.NullBool{Bool: false, Valid: true},
		Warnings:        []string{"42:this is a warning"},
		Data:            testPipelineData(),
	}

	p := new(api.Pipeline)

	p.SetID(1)
	p.SetRepo(testRepo().ToAPI())
	p.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	p.SetFlavor("large")
	p.SetPlatform("docker")
	p.SetRef("refs/heads/main")
	p.SetType(constants.PipelineTypeYAML)
	p.SetVersion("1")
	p.SetExternalSecrets(false)
	p.SetInternalSecrets(false)
	p.SetServices(true)
	p.SetStages(false)
	p.SetSteps(true)
	p.SetTemplates(false)
	p.SetArtifact(false)
	p.SetWarnings([]string{"42:this is a warning"})
	p.SetData(testPipelineData())

	// run test
	got := PipelineFromAPI(p)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("PipelineFromAPI is %v, want %v", got, want)
	}
}

// testPipeline is a test helper function to create a Pipeline
// type with all fields set to a fake value.
func testPipeline() *Pipeline {
	return &Pipeline{
		ID:              sql.NullInt64{Int64: 1, Valid: true},
		RepoID:          sql.NullInt64{Int64: 1, Valid: true},
		Commit:          sql.NullString{String: "48afb5bdc41ad69bf22588491333f7cf71135163", Valid: true},
		Flavor:          sql.NullString{String: "large", Valid: true},
		Platform:        sql.NullString{String: "docker", Valid: true},
		Ref:             sql.NullString{String: "refs/heads/main", Valid: true},
		Type:            sql.NullString{String: constants.PipelineTypeYAML, Valid: true},
		Version:         sql.NullString{String: "1", Valid: true},
		ExternalSecrets: sql.NullBool{Bool: false, Valid: true},
		InternalSecrets: sql.NullBool{Bool: false, Valid: true},
		Services:        sql.NullBool{Bool: true, Valid: true},
		Stages:          sql.NullBool{Bool: false, Valid: true},
		Steps:           sql.NullBool{Bool: true, Valid: true},
		Templates:       sql.NullBool{Bool: false, Valid: true},
		Artifacts:       sql.NullBool{Bool: false, Valid: true},
		Warnings:        []string{"42:this is a warning"},
		Data:            testPipelineData(),

		Repo: *testRepo(),
	}
}

// testPipelineData is a test helper function to create the
// content for the Data field for the Pipeline type.
func testPipelineData() []byte {
	return []byte(`
version: 1

worker:
  flavor: large
  platform: docker

services:
  - name: redis
    image: redis

steps:
  - name: ping
    image: redis
    commands:
      - redis-cli -h redis ping
`)
}
