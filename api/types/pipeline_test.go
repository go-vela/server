// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAPI_Pipeline_Getters(t *testing.T) {
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
			pipeline: new(Pipeline),
			want:     new(Pipeline),
		},
	}

	// run tests
	for _, test := range tests {
		if test.pipeline.GetID() != test.want.GetID() {
			t.Errorf("GetID is %v, want %v", test.pipeline.GetID(), test.want.GetID())
		}

		if !reflect.DeepEqual(test.pipeline.GetRepo(), test.want.GetRepo()) {
			t.Errorf("GetRepoID is %v, want %v", test.pipeline.GetRepo(), test.want.GetRepo())
		}

		if test.pipeline.GetCommit() != test.want.GetCommit() {
			t.Errorf("GetCommit is %v, want %v", test.pipeline.GetCommit(), test.want.GetCommit())
		}

		if test.pipeline.GetFlavor() != test.want.GetFlavor() {
			t.Errorf("GetFlavor is %v, want %v", test.pipeline.GetFlavor(), test.want.GetFlavor())
		}

		if test.pipeline.GetPlatform() != test.want.GetPlatform() {
			t.Errorf("GetPlatform is %v, want %v", test.pipeline.GetPlatform(), test.want.GetPlatform())
		}

		if test.pipeline.GetRef() != test.want.GetRef() {
			t.Errorf("GetRef is %v, want %v", test.pipeline.GetRef(), test.want.GetRef())
		}

		if test.pipeline.GetType() != test.want.GetType() {
			t.Errorf("GetType is %v, want %v", test.pipeline.GetType(), test.want.GetType())
		}

		if test.pipeline.GetVersion() != test.want.GetVersion() {
			t.Errorf("GetVersion is %v, want %v", test.pipeline.GetVersion(), test.want.GetVersion())
		}

		if test.pipeline.GetExternalSecrets() != test.want.GetExternalSecrets() {
			t.Errorf("GetExternalSecrets is %v, want %v", test.pipeline.GetExternalSecrets(), test.want.GetExternalSecrets())
		}

		if test.pipeline.GetInternalSecrets() != test.want.GetInternalSecrets() {
			t.Errorf("GetInternalSecrets is %v, want %v", test.pipeline.GetInternalSecrets(), test.want.GetInternalSecrets())
		}

		if test.pipeline.GetServices() != test.want.GetServices() {
			t.Errorf("GetServices is %v, want %v", test.pipeline.GetServices(), test.want.GetServices())
		}

		if test.pipeline.GetStages() != test.want.GetStages() {
			t.Errorf("GetStages is %v, want %v", test.pipeline.GetStages(), test.want.GetStages())
		}

		if test.pipeline.GetSteps() != test.want.GetSteps() {
			t.Errorf("GetSteps is %v, want %v", test.pipeline.GetSteps(), test.want.GetSteps())
		}

		if test.pipeline.GetTemplates() != test.want.GetTemplates() {
			t.Errorf("GetTemplates is %v, want %v", test.pipeline.GetTemplates(), test.want.GetTemplates())
		}

		if !reflect.DeepEqual(test.pipeline.GetWarnings(), test.want.GetWarnings()) {
			t.Errorf("GetWarnings is %v, want %v", test.pipeline.GetWarnings(), test.want.GetWarnings())
		}

		if !reflect.DeepEqual(test.pipeline.GetData(), test.want.GetData()) {
			t.Errorf("GetData is %v, want %v", test.pipeline.GetData(), test.want.GetData())
		}
	}
}

func TestAPI_Pipeline_Setters(t *testing.T) {
	// setup types
	var p *Pipeline

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
			want:     new(Pipeline),
		},
	}

	// run tests
	for _, test := range tests {
		test.pipeline.SetID(test.want.GetID())
		test.pipeline.SetRepo(test.want.GetRepo())
		test.pipeline.SetCommit(test.want.GetCommit())
		test.pipeline.SetFlavor(test.want.GetFlavor())
		test.pipeline.SetPlatform(test.want.GetPlatform())
		test.pipeline.SetRef(test.want.GetRef())
		test.pipeline.SetType(test.want.GetType())
		test.pipeline.SetVersion(test.want.GetVersion())
		test.pipeline.SetExternalSecrets(test.want.GetExternalSecrets())
		test.pipeline.SetInternalSecrets(test.want.GetInternalSecrets())
		test.pipeline.SetServices(test.want.GetServices())
		test.pipeline.SetStages(test.want.GetStages())
		test.pipeline.SetSteps(test.want.GetSteps())
		test.pipeline.SetTemplates(test.want.GetTemplates())
		test.pipeline.SetWarnings(test.want.GetWarnings())
		test.pipeline.SetData(test.want.GetData())

		if test.pipeline.GetID() != test.want.GetID() {
			t.Errorf("SetID is %v, want %v", test.pipeline.GetID(), test.want.GetID())
		}

		if !reflect.DeepEqual(test.pipeline.GetRepo(), test.want.GetRepo()) {
			t.Errorf("SetRepoID is %v, want %v", test.pipeline.GetRepo(), test.want.GetRepo())
		}

		if test.pipeline.GetCommit() != test.want.GetCommit() {
			t.Errorf("SetCommit is %v, want %v", test.pipeline.GetCommit(), test.want.GetCommit())
		}

		if test.pipeline.GetFlavor() != test.want.GetFlavor() {
			t.Errorf("SetFlavor is %v, want %v", test.pipeline.GetFlavor(), test.want.GetFlavor())
		}

		if test.pipeline.GetPlatform() != test.want.GetPlatform() {
			t.Errorf("SetPlatform is %v, want %v", test.pipeline.GetPlatform(), test.want.GetPlatform())
		}

		if test.pipeline.GetRef() != test.want.GetRef() {
			t.Errorf("SetRef is %v, want %v", test.pipeline.GetRef(), test.want.GetRef())
		}

		if test.pipeline.GetType() != test.want.GetType() {
			t.Errorf("SetType is %v, want %v", test.pipeline.GetType(), test.want.GetType())
		}

		if test.pipeline.GetVersion() != test.want.GetVersion() {
			t.Errorf("SetVersion is %v, want %v", test.pipeline.GetVersion(), test.want.GetVersion())
		}

		if test.pipeline.GetExternalSecrets() != test.want.GetExternalSecrets() {
			t.Errorf("SetExternalSecrets is %v, want %v", test.pipeline.GetExternalSecrets(), test.want.GetExternalSecrets())
		}

		if test.pipeline.GetInternalSecrets() != test.want.GetInternalSecrets() {
			t.Errorf("SetInternalSecrets is %v, want %v", test.pipeline.GetInternalSecrets(), test.want.GetInternalSecrets())
		}

		if test.pipeline.GetServices() != test.want.GetServices() {
			t.Errorf("SetServices is %v, want %v", test.pipeline.GetServices(), test.want.GetServices())
		}

		if test.pipeline.GetStages() != test.want.GetStages() {
			t.Errorf("SetStages is %v, want %v", test.pipeline.GetStages(), test.want.GetStages())
		}

		if test.pipeline.GetSteps() != test.want.GetSteps() {
			t.Errorf("SetSteps is %v, want %v", test.pipeline.GetSteps(), test.want.GetSteps())
		}

		if test.pipeline.GetTemplates() != test.want.GetTemplates() {
			t.Errorf("SetTemplates is %v, want %v", test.pipeline.GetTemplates(), test.want.GetTemplates())
		}

		if !reflect.DeepEqual(test.pipeline.GetWarnings(), test.want.GetWarnings()) {
			t.Errorf("SetWarnings is %v, want %v", test.pipeline.GetWarnings(), test.want.GetWarnings())
		}

		if !reflect.DeepEqual(test.pipeline.GetData(), test.want.GetData()) {
			t.Errorf("SetData is %v, want %v", test.pipeline.GetData(), test.want.GetData())
		}
	}
}

func TestAPI_Pipeline_String(t *testing.T) {
	// setup types
	p := testPipeline()

	want := fmt.Sprintf(`{
  Commit: %s,
  Data: %s,
  Flavor: %s,
  ID: %d,
  Platform: %s,
  Ref: %s,
  Repo: %v,
  ExternalSecrets: %t,
  InternalSecrets: %t,
  Services: %t,
  Stages: %t,
  Steps: %t,
  Templates: %t,
  TestReport: %t,
  Type: %s,
  Version: %s,
  Warnings: %v,
}`,
		p.GetCommit(),
		p.GetData(),
		p.GetFlavor(),
		p.GetID(),
		p.GetPlatform(),
		p.GetRef(),
		p.GetRepo(),
		p.GetExternalSecrets(),
		p.GetInternalSecrets(),
		p.GetServices(),
		p.GetStages(),
		p.GetSteps(),
		p.GetTemplates(),
		p.GetTestReport(),
		p.GetType(),
		p.GetVersion(),
		p.GetWarnings(),
	)

	// run test
	got := p.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

// testPipeline is a test helper function to create a Pipeline
// type with all fields set to a fake value.
func testPipeline() *Pipeline {
	p := new(Pipeline)

	p.SetID(1)
	p.SetRepo(testRepo())
	p.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	p.SetFlavor("large")
	p.SetPlatform("docker")
	p.SetRef("refs/heads/main")
	p.SetRef("yaml")
	p.SetVersion("1")
	p.SetExternalSecrets(false)
	p.SetInternalSecrets(false)
	p.SetServices(true)
	p.SetStages(false)
	p.SetSteps(true)
	p.SetTemplates(false)
	p.SetTestReport(true)
	p.SetData(testPipelineData())
	p.SetWarnings([]string{"42:this is a warning"})

	return p
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
