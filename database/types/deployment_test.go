// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lib/pq"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/raw"
)

func TestDatabase_Deployment_Nullify(t *testing.T) {
	// setup types
	var d *Deployment

	want := &Deployment{
		ID:          sql.NullInt64{Int64: 0, Valid: false},
		Number:      sql.NullInt64{Int64: 0, Valid: false},
		RepoID:      sql.NullInt64{Int64: 0, Valid: false},
		URL:         sql.NullString{String: "", Valid: false},
		Commit:      sql.NullString{String: "", Valid: false},
		Ref:         sql.NullString{String: "", Valid: false},
		Task:        sql.NullString{String: "", Valid: false},
		Target:      sql.NullString{String: "", Valid: false},
		Description: sql.NullString{String: "", Valid: false},
		Payload:     nil,
		CreatedAt:   sql.NullInt64{Int64: 0, Valid: false},
		CreatedBy:   sql.NullString{String: "", Valid: false},
		Builds:      nil,
	}

	// setup tests
	tests := []struct {
		deployment *Deployment
		want       *Deployment
	}{
		{
			deployment: testDeployment(),
			want:       testDeployment(),
		},
		{
			deployment: d,
			want:       nil,
		},
		{
			deployment: new(Deployment),
			want:       want,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.deployment.Nullify()

		if diff := cmp.Diff(test.want, got); diff != "" {
			t.Errorf("(ToAPI: -want +got):\n%s", diff)
		}
	}
}

func TestDatabase_Deployment_ToAPI(t *testing.T) {
	builds := []*api.Build{testBuild().ToAPI()}

	want := new(api.Deployment)
	want.SetID(1)
	want.SetNumber(1)
	want.SetRepo(testRepo().ToAPI())
	want.SetURL("https://github.com/github/octocat/deployments/1")
	want.SetCommit("1234")
	want.SetRef("refs/heads/main")
	want.SetTask("deploy:vela")
	want.SetTarget("production")
	want.SetDescription("Deployment request from Vela")
	want.SetPayload(raw.StringSliceMap{"foo": "test1"})
	want.SetCreatedAt(1)
	want.SetCreatedBy("octocat")
	want.SetBuilds(builds)

	got := testDeployment().ToAPI(builds)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("(ToAPI: -want +got):\n%s", diff)
	}
}

func TestDatabase_Deployment_Validate(t *testing.T) {
	// setup types
	tests := []struct {
		failure    bool
		deployment *Deployment
		want       *Deployment
	}{
		{
			failure:    false,
			deployment: testDeployment(),
		},
		{ // no number set for deployment
			failure: true,
			deployment: &Deployment{
				ID:     sql.NullInt64{Int64: 1, Valid: true},
				RepoID: sql.NullInt64{Int64: 1, Valid: true},
			},
			want: &Deployment{
				ID:     sql.NullInt64{Int64: 1, Valid: true},
				RepoID: sql.NullInt64{Int64: 1, Valid: true},
			},
		},
		{ // no repoID set for deployment
			failure: true,
			deployment: &Deployment{
				ID:     sql.NullInt64{Int64: 1, Valid: true},
				Number: sql.NullInt64{Int64: 1, Valid: true},
			},
			want: &Deployment{
				ID:     sql.NullInt64{Int64: 1, Valid: true},
				RepoID: sql.NullInt64{Int64: 1, Valid: true},
			},
		},
		{ // too many builds
			failure: true,
			deployment: &Deployment{
				ID:     sql.NullInt64{Int64: 1, Valid: true},
				Number: sql.NullInt64{Int64: 1, Valid: true},
				Builds: generateBuilds(100),
			},
			want: &Deployment{
				ID:     sql.NullInt64{Int64: 1, Valid: true},
				RepoID: sql.NullInt64{Int64: 1, Valid: true},
				Builds: generateBuilds(50),
			},
		},
		{ // acceptable builds
			failure: true,
			deployment: &Deployment{
				ID:     sql.NullInt64{Int64: 1, Valid: true},
				Number: sql.NullInt64{Int64: 1, Valid: true},
				Builds: generateBuilds(30),
			},
			want: &Deployment{
				ID:     sql.NullInt64{Int64: 1, Valid: true},
				RepoID: sql.NullInt64{Int64: 1, Valid: true},
				Builds: generateBuilds(30),
			},
		},
	}

	// run tests
	for _, test := range tests {
		err := test.deployment.Validate()

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

func TestDatabase_DeploymentFromAPI(t *testing.T) {
	builds := []*api.Build{testBuild().ToAPI()}

	d := new(api.Deployment)
	d.SetID(1)
	d.SetNumber(1)
	d.SetRepo(testRepo().ToAPI())
	d.SetURL("https://github.com/github/octocat/deployments/1")
	d.SetCommit("1234")
	d.SetRef("refs/heads/main")
	d.SetTask("deploy:vela")
	d.SetTarget("production")
	d.SetDescription("Deployment request from Vela")
	d.SetPayload(raw.StringSliceMap{"foo": "test1"})
	d.SetCreatedAt(1)
	d.SetCreatedBy("octocat")
	d.SetBuilds(builds)

	want := &Deployment{
		ID:          sql.NullInt64{Int64: 1, Valid: true},
		Number:      sql.NullInt64{Int64: 1, Valid: true},
		RepoID:      sql.NullInt64{Int64: 1, Valid: true},
		URL:         sql.NullString{String: "https://github.com/github/octocat/deployments/1", Valid: true},
		Commit:      sql.NullString{String: "1234", Valid: true},
		Ref:         sql.NullString{String: "refs/heads/main", Valid: true},
		Task:        sql.NullString{String: "deploy:vela", Valid: true},
		Target:      sql.NullString{String: "production", Valid: true},
		Description: sql.NullString{String: "Deployment request from Vela", Valid: true},
		Payload:     raw.StringSliceMap{"foo": "test1"},
		CreatedAt:   sql.NullInt64{Int64: 1, Valid: true},
		CreatedBy:   sql.NullString{String: "octocat", Valid: true},
		Builds:      pq.StringArray{"1"},
	}

	// run test
	got := DeploymentFromAPI(d)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
}

// testDeployment is a test helper function to create a Deployment type with all fields set to a fake value.
func testDeployment() *Deployment {
	return &Deployment{
		ID:          sql.NullInt64{Int64: 1, Valid: true},
		Number:      sql.NullInt64{Int64: 1, Valid: true},
		RepoID:      sql.NullInt64{Int64: 1, Valid: true},
		URL:         sql.NullString{String: "https://github.com/github/octocat/deployments/1", Valid: true},
		Commit:      sql.NullString{String: "1234", Valid: true},
		Ref:         sql.NullString{String: "refs/heads/main", Valid: true},
		Task:        sql.NullString{String: "deploy:vela", Valid: true},
		Target:      sql.NullString{String: "production", Valid: true},
		Description: sql.NullString{String: "Deployment request from Vela", Valid: true},
		Payload:     raw.StringSliceMap{"foo": "test1"},
		CreatedAt:   sql.NullInt64{Int64: 1, Valid: true},
		CreatedBy:   sql.NullString{String: "octocat", Valid: true},
		Builds:      pq.StringArray{"1"},

		Repo: *testRepo(),
	}
}

// generateBuilds returns a list of valid builds that exceed the maximum size.
func generateBuilds(amount int) []string {
	// initialize empty builds
	builds := []string{}

	for range amount {
		builds = append(builds, "123456789")
	}

	return builds
}
