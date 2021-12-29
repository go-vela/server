// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"
	"testing"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

// TODO: update the test client to not return nil
// https://github.com/jenkins-x/go-scm/blob/main/scm/driver/fake/content.go#L43
// func TestNative_ConfigBackoff_YML(t *testing.T) {
// 	// setup types
// 	u := new(library.User)
// 	u.SetName("foo")
// 	u.SetToken("bar")

// 	r := new(library.Repo)
// 	r.SetOrg("foo")
// 	r.SetName("bar")

// 	want, err := ioutil.ReadFile("testdata/pipeline.yml")
// 	if err != nil {
// 		// t.Errorf("Config reading file returned err: %v", err)
// 	}

// 	client, _ := NewTest("fake.com")

// 	// run test
// 	got, err := client.Config(u, r, "")

// 	if err != nil {
// 		t.Errorf("Config returned err: %v", err)
// 	}

// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("Config is %v, want %v", got, want)
// 	}
// }

// TODO: update the test client to not return nil
// https://github.com/jenkins-x/go-scm/blob/main/scm/driver/fake/content.go#L43
// func TestNative_Config_YML(t *testing.T) {
// 	// setup types
// 	u := new(library.User)
// 	u.SetName("foo")
// 	u.SetToken("bar")

// 	r := new(library.Repo)
// 	r.SetOrg("foo")
// 	r.SetName("bar")

// 	want, err := ioutil.ReadFile("testdata/pipeline.yml")
// 	if err != nil {
// 		// t.Errorf("Config reading file returned err: %v", err)
// 	}

// 	client, _ := NewTest("fake.com")

// 	// run test
// 	got, err := client.Config(u, r, "")

// 	if err != nil {
// 		t.Errorf("Config returned err: %v", err)
// 	}

// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("Config is %v, want %v", got, want)
// 	}
// }

func TestGithub_Disable(t *testing.T) {
	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest("fake.com")

	// run test
	err := client.Disable(u, "foo", "bar")

	if err != nil {
		t.Errorf("Disable returned err: %v", err)
	}
}

// TODO: update the test client to not return nil
// https://github.com/jenkins-x/go-scm/blob/main/scm/driver/fake/repo.go#L192
// func TestGithub_Enable(t *testing.T) {
// 	// setup types
// 	u := new(library.User)
// 	u.SetName("foo")
// 	u.SetToken("bar")

// 	client, _ := NewTest("fake.com")

// 	// run test
// 	_, err := client.Enable(u, "foo", "bar", "secret")

// 	if err != nil {
// 		t.Errorf("Enable returned err: %v", err)
// 	}
// }

func TestGithub_Status_Deployment(t *testing.T) {
	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	b := new(library.Build)
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)
	b.SetEvent(constants.EventDeploy)
	b.SetStatus(constants.StatusRunning)
	b.SetCommit("abcd1234")
	b.SetSource(fmt.Sprintf("%s/%s/%s/deployments/1", "fake.com", "foo", "bar"))

	client, _ := NewTest("fake.com")

	// run test
	err := client.Status(u, b, "foo", "bar")

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}
}
