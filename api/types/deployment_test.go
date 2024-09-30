// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTypes_Deployment_Getters(t *testing.T) {
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
			deployment: new(Deployment),
			want:       new(Deployment),
		},
	}

	// run tests
	for _, test := range tests {
		if test.deployment.GetID() != test.want.GetID() {
			t.Errorf("GetID is %v, want %v", test.deployment.GetID(), test.want.GetID())
		}

		if !reflect.DeepEqual(test.deployment.GetRepo(), test.want.GetRepo()) {
			t.Errorf("GetRepoID is %v, want %v", test.deployment.GetRepo(), test.want.GetRepo())
		}

		if test.deployment.GetURL() != test.want.GetURL() {
			t.Errorf("GetURL is %v, want %v", test.deployment.GetURL(), test.want.GetURL())
		}

		if test.deployment.GetCommit() != test.want.GetCommit() {
			t.Errorf("GetCommit is %v, want %v", test.deployment.GetCommit(), test.want.GetCommit())
		}

		if test.deployment.GetRef() != test.want.GetRef() {
			t.Errorf("GetRef is %v, want %v", test.deployment.GetRef(), test.want.GetRef())
		}

		if test.deployment.GetTask() != test.want.GetTask() {
			t.Errorf("GetTask is %v, want %v", test.deployment.GetTask(), test.want.GetTask())
		}

		if test.deployment.GetTarget() != test.want.GetTarget() {
			t.Errorf("GetTarget is %v, want %v", test.deployment.GetTarget(), test.want.GetTarget())
		}

		if test.deployment.GetDescription() != test.want.GetDescription() {
			t.Errorf("GetDescription is %v, want %v", test.deployment.GetDescription(), test.want.GetDescription())
		}

		if !reflect.DeepEqual(test.deployment.GetPayload(), test.want.GetPayload()) {
			t.Errorf("GetPayload is %v, want %v", test.deployment.GetPayload(), test.want.GetPayload())
		}

		if test.deployment.GetCreatedAt() != test.want.GetCreatedAt() {
			t.Errorf("GetCreatedAt is %v, want %v", test.deployment.GetCreatedAt(), test.want.GetCreatedAt())
		}

		if test.deployment.GetCreatedBy() != test.want.GetCreatedBy() {
			t.Errorf("GetCreatedBy is %v, want %v", test.deployment.GetCreatedBy(), test.want.GetCreatedBy())
		}
	}
}

func TestTypes_Deployment_Setters(t *testing.T) {
	// setup types
	var d *Deployment

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
			want:       new(Deployment),
		},
	}

	// run tests
	for _, test := range tests {
		test.deployment.SetID(test.want.GetID())
		test.deployment.SetNumber(test.want.GetNumber())
		test.deployment.SetRepo(test.want.GetRepo())
		test.deployment.SetURL(test.want.GetURL())
		test.deployment.SetCommit(test.want.GetCommit())
		test.deployment.SetRef(test.want.GetRef())
		test.deployment.SetTask(test.want.GetTask())
		test.deployment.SetTarget(test.want.GetTarget())
		test.deployment.SetDescription(test.want.GetDescription())
		test.deployment.SetPayload(test.want.GetPayload())
		test.deployment.SetCreatedAt(test.want.GetCreatedAt())
		test.deployment.SetCreatedBy(test.want.GetCreatedBy())

		if test.deployment.GetID() != test.want.GetID() {
			t.Errorf("SetID is %v, want %v", test.deployment.GetID(), test.want.GetID())
		}

		if test.deployment.GetNumber() != test.want.GetNumber() {
			t.Errorf("SetNumber is %v, want %v", test.deployment.GetNumber(), test.want.GetNumber())
		}

		if !reflect.DeepEqual(test.deployment.GetRepo(), test.want.GetRepo()) {
			t.Errorf("SetRepoID is %v, want %v", test.deployment.GetRepo(), test.want.GetRepo())
		}

		if test.deployment.GetURL() != test.want.GetURL() {
			t.Errorf("SetURL is %v, want %v", test.deployment.GetURL(), test.want.GetURL())
		}

		if test.deployment.GetCommit() != test.want.GetCommit() {
			t.Errorf("SetCommit is %v, want %v", test.deployment.GetCommit(), test.want.GetCommit())
		}

		if test.deployment.GetRef() != test.want.GetRef() {
			t.Errorf("SetRef is %v, want %v", test.deployment.GetRef(), test.want.GetRef())
		}

		if test.deployment.GetTask() != test.want.GetTask() {
			t.Errorf("SetTask is %v, want %v", test.deployment.GetTask(), test.want.GetTask())
		}

		if test.deployment.GetTarget() != test.want.GetTarget() {
			t.Errorf("SetTarget is %v, want %v", test.deployment.GetTarget(), test.want.GetTarget())
		}

		if test.deployment.GetDescription() != test.want.GetDescription() {
			t.Errorf("SetDescription is %v, want %v", test.deployment.GetDescription(), test.want.GetDescription())
		}

		if !reflect.DeepEqual(test.deployment.GetPayload(), test.want.GetPayload()) {
			t.Errorf("SetPayload is %v, want %v", test.deployment.GetPayload(), test.want.GetPayload())
		}

		if test.deployment.GetCreatedAt() != test.want.GetCreatedAt() {
			t.Errorf("SetCreatedAt is %v, want %v", test.deployment.GetCreatedAt(), test.want.GetCreatedAt())
		}

		if test.deployment.GetCreatedBy() != test.want.GetCreatedBy() {
			t.Errorf("SetCreatedBy is %v, want %v", test.deployment.GetCreatedBy(), test.want.GetCreatedBy())
		}
	}
}

func TestTypes_Deployment_String(t *testing.T) {
	// setup types
	d := testDeployment()

	want := fmt.Sprintf(`{
  Commit: %s,
  CreatedAt: %d,
  CreatedBy: %s,
  Description: %s,
  ID: %d,
  Number: %d,
  Ref: %s,
  Repo: %v,
  Target: %s,
  Task: %s,
  URL: %s,
  Payload: %s,
  Builds: %d,
}`,
		d.GetCommit(),
		d.GetCreatedAt(),
		d.GetCreatedBy(),
		d.GetDescription(),
		d.GetID(),
		d.GetNumber(),
		d.GetRef(),
		d.GetRepo(),
		d.GetTarget(),
		d.GetTask(),
		d.GetURL(),
		d.GetPayload(),
		len(d.GetBuilds()),
	)

	// run test
	got := d.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

// testDeployment is a test helper function to create a Deployment
// type with all fields set to a fake value.
func testDeployment() *Deployment {
	d := new(Deployment)

	d.SetID(1)
	d.SetNumber(0)
	d.SetRepo(testRepo())
	d.SetURL("https://api.github.com/repos/github/octocat/deployments/1")
	d.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	d.SetRef("refs/heads/main")
	d.SetTask("vela-deploy")
	d.SetTarget("production")
	d.SetDescription("Deployment request from Vela")
	d.SetPayload(map[string]string{
		"foo": "test1",
	})
	d.SetCreatedAt(1)
	d.SetCreatedBy("octocat")

	return d
}
