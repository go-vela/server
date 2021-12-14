package native

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/raw"
)

func TestNative_GetDeployment(t *testing.T) {
	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(library.Repo)
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	want := new(library.Deployment)
	want.SetID(1)
	want.SetRepoID(1)
	want.SetURL("https://api.github.com/repos/foo/bar/deployments/1")
	want.SetUser("octocat")
	want.SetCommit("a84d88e7554fc1fa21bcbc4efae3c782a70d2b9d")
	want.SetRef("topic-branch")
	want.SetTask("deploy")
	want.SetTarget("production")
	want.SetDescription("Deploy request from Vela")
	want.SetPayload(raw.StringSliceMap{"deploy": "migrate"})

	client, _ := NewTest("fake.com")

	// run test
	got, err := client.GetDeployment(u, r, 1)

	if err != nil {
		t.Errorf("GetDeployment returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetDeployment is %v, want %v", got, want)
	}
}

// TODO: update the test client to not return nil
// https://github.com/jenkins-x/go-scm/blob/main/scm/driver/fake/deploy.go#L26
// func TestNative_GetDeploymentCount(t *testing.T) {
// 	// setup types
// 	u := new(library.User)
// 	u.SetName("foo")
// 	u.SetToken("bar")

// 	r := new(library.Repo)
// 	r.SetID(1)
// 	r.SetOrg("foo")
// 	r.SetName("bar")
// 	r.SetFullName("foo/bar")

// 	want := int64(2)

// 	client, _ := NewTest("fake.com")

// 	// run test
// 	got, err := client.GetDeploymentCount(u, r)

// 	if err != nil {
// 		t.Errorf("GetDeployment returned err: %v", err)
// 	}

// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("GetDeployment is %v, want %v", got, want)
// 	}
// }

func TestGithub_GetDeploymentList(t *testing.T) {
	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(library.Repo)
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	d1 := new(library.Deployment)
	d1.SetID(1)
	d1.SetRepoID(1)
	d1.SetURL("https://api.github.com/repos/foo/bar/deployments/1")
	d1.SetUser("octocat")
	d1.SetCommit("a84d88e7554fc1fa21bcbc4efae3c782a70d2b9d")
	d1.SetRef("topic-branch")
	d1.SetTask("deploy")
	d1.SetTarget("production")
	d1.SetDescription("Deploy request from Vela")
	d1.SetPayload(raw.StringSliceMap{"deploy": "migrate"})

	want := []*library.Deployment{d1}

	client, _ := NewTest("fake.com")

	// run test
	got, err := client.GetDeploymentList(u, r, 1, 100)

	if err != nil {
		t.Errorf("GetDeployment returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetDeployment is %v, want %v", got, want)
	}
}

// TODO: update the test client have a valid author
// https://github.com/jenkins-x/go-scm/blob/main/scm/driver/fake/deploy.go#L47
// func TestGithub_CreateDeployment(t *testing.T) {
// 	// setup types
// 	u := new(library.User)
// 	u.SetName("foo")
// 	u.SetToken("bar")

// 	r := new(library.Repo)
// 	r.SetID(1)
// 	r.SetOrg("foo")
// 	r.SetName("bar")
// 	r.SetFullName("foo/bar")

// 	d := new(library.Deployment)
// 	d.SetID(1)
// 	d.SetRepoID(1)
// 	d.SetURL("https://api.github.com/repos/foo/bar/deployments/1")
// 	d.SetUser("octocat")
// 	d.SetCommit("a84d88e7554fc1fa21bcbc4efae3c782a70d2b9d")
// 	d.SetRef("topic-branch")
// 	d.SetTask("deploy")
// 	d.SetTarget("production")
// 	d.SetDescription("Deploy request from Vela")

// 	client, _ := NewTest("fake.com")

// 	// run test
// 	err := client.CreateDeployment(u, r, d)

// 	if err != nil {
// 		t.Errorf("CreateDeployment returned err: %v", err)
// 	}
// }
