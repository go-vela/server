// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
)

func TestNative_OrgAccess_Admin(t *testing.T) {
	// setup types
	want := "admin"

	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest("fake.com")

	// run test
	got, err := client.OrgAccess(u, "test-org1")

	if err != nil {
		t.Errorf("OrgAccess returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("OrgAccess is %v, want %v", got, want)
	}
}

// TODO: Add support in Jenkins-X for custom data with Membership
// https://github.com/jenkins-x/go-scm/blob/main/scm/driver/fake/org.go#L120

// func TestNative_OrgAccess_Member(t *testing.T) {
// 	// setup types
// 	want := "member"

// 	u := new(library.User)
// 	u.SetName("foo")
// 	u.SetToken("bar")

// 	client, _ := NewTest("fake.com")

// 	// run test
// 	got, err := client.OrgAccess(u, "test-org2")

// 	if err != nil {
// 		t.Errorf("OrgAccess returned err: %v", err)
// 	}

// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("OrgAccess is %v, want %v", got, want)
// 	}
// }

// TODO: Add support in Jenkins-X for custom data with Membership
// https://github.com/jenkins-x/go-scm/blob/main/scm/driver/fake/org.go#L120
//
// func TestNative_OrgAccess_NotFound(t *testing.T) {
// setup types
// 	want := ""

// 	u := new(library.User)
// 	u.SetName("foo")
// 	u.SetToken("bar")

// 	client, _ := NewTest("fake.com")

// 	// run test
// 	got, err := client.OrgAccess(u, "github")

// 	if err == nil {
// 		t.Errorf("OrgAccess should have returned err")
// 	}

// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("OrgAccess is %v, want %v", got, want)
// 	}
// }

// TODO: Add something to make pagination less painful in testing
// https://github.com/jenkins-x/go-scm/blob/main/scm/driver/fake/org.go#L133
// func TestNative_OrgAccess_Pending(t *testing.T) {
// 	// setup types
// 	want := ""

// 	u := new(library.User)
// 	u.SetName("foo")
// 	u.SetToken("bar")

// 	client, _ := NewTest("fake.com")

// 	// run test
// 	got, err := client.OrgAccess(u, "test-org2")

// 	if err != nil {
// 		t.Errorf("OrgAccess returned err: %v", err)
// 	}

// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("OrgAccess is %v, want %v", got, want)
// 	}
// }

// TODO: Add support in Jenkins-X for custom data with Membership
// https://github.com/jenkins-x/go-scm/blob/main/scm/driver/fake/org.go#L120
// func TestNative_OrgAccess_Personal(t *testing.T) {
// 	// setup types
// 	want := "admin"

// 	u := new(library.User)
// 	u.SetName("foo")
// 	u.SetToken("bar")

// 	client, _ := NewTest("fake.com")

// 	// run test
// 	got, err := client.OrgAccess(u, "foo")

// 	if err != nil {
// 		t.Errorf("OrgAccess returned err: %v", err)
// 	}

// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("OrgAccess is %v, want %v", got, want)
// 	}
// }

func TestNative_RepoAccess_Admin(t *testing.T) {
	// setup types
	want := "admin"

	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest("fake.com")

	// run test
	got, err := client.RepoAccess(u, u.GetToken(), "github", "octocat")

	if err != nil {
		t.Errorf("RepoAccess returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("RepoAccess is %v, want %v", got, want)
	}
}

func TestNative_RepoAccess_NotFound(t *testing.T) {
	// setup types
	want := ""

	u := new(library.User)
	u.SetName("notfound")
	u.SetToken("bar")

	client, _ := NewTest("fake.com")

	// run test
	got, _ := client.RepoAccess(u, u.GetToken(), "github", "octocat")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("RepoAccess is %v, want %v", got, want)
	}
}

// TODO: Update fake ListTeams to work with pagination
// https://github.com/jenkins-x/go-scm/blob/main/scm/driver/fake/org.go#L86
// func TestNative_TeamAccess_Admin(t *testing.T) {
// 	// setup types
// 	want := "admin"

// 	u := new(library.User)
// 	u.SetName("foo")
// 	u.SetToken("bar")

// 	client, _ := NewTest("fake.com")

// 	// run test
// 	got, err := client.TeamAccess(u, "github", "octocat")

// 	if err != nil {
// 		t.Errorf("TeamAccess returned err: %v", err)
// 	}

// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("TeamAccess is %v, want %v", got, want)
// 	}
// }

// TODO: Update fake ListTeams to work with pagination
// https://github.com/jenkins-x/go-scm/blob/main/scm/driver/fake/org.go#L86
// func TestNative_TeamAccess_NoAccess(t *testing.T) {
// 	// setup context
// 	gin.SetMode(gin.TestMode)

// 	resp := httptest.NewRecorder()
// 	_, engine := gin.CreateTestContext(resp)

// 	// setup mock server
// 	engine.GET("/api/v3/user/teams", func(c *gin.Context) {
// 		c.Header("Content-Type", "application/json")
// 		c.Status(http.StatusOK)
// 		c.File("testdata/team_admin.json")
// 	})

// 	s := httptest.NewServer(engine)
// 	defer s.Close()

// 	// setup types
// 	want := ""

// 	u := new(library.User)
// 	u.SetName("foo")
// 	u.SetToken("bar")

// 	client, _ := NewTest(s.URL)

// 	// run test
// 	got, err := client.TeamAccess(u, "github", "baz")

// 	if resp.Code != http.StatusOK {
// 		t.Errorf("TeamAccess returned %v, want %v", resp.Code, http.StatusOK)
// 	}

// 	if err != nil {
// 		t.Errorf("TeamAccess returned err: %v", err)
// 	}

// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("TeamAccess is %v, want %v", got, want)
// 	}
// }

// TODO: Update fake ListTeams to work with pagination
// https://github.com/jenkins-x/go-scm/blob/main/scm/driver/fake/org.go#L86
// func TestNative_TeamAccess_NotFound(t *testing.T) {
// 	// setup mock server
// 	s := httptest.NewServer(http.NotFoundHandler())
// 	defer s.Close()

// 	// setup types
// 	want := ""

// 	u := new(library.User)
// 	u.SetName("foo")
// 	u.SetToken("bar")

// 	client, _ := NewTest(s.URL)

// 	// run test
// 	got, err := client.TeamAccess(u, "github", "octocat")

// 	if err == nil {
// 		t.Errorf("TeamAccess should have returned err")
// 	}

// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("TeamAccess is %v, want %v", got, want)
// 	}
// }

// TODO: Update fake ListTeams to work with pagination
// https://github.com/jenkins-x/go-scm/blob/main/scm/driver/fake/org.go#L86
// func TestGithub_TeamList(t *testing.T) {
// 	// setup types
// 	want := []string{"Justice League", "octocat"}

// 	u := new(library.User)
// 	u.SetName("foo")
// 	u.SetToken("bar")

// 	client, _ := NewTest("fake.com")

// 	// run test
// 	got, err := client.ListUsersTeamsForOrg(u, "github")

// 	if err != nil {
// 		t.Errorf("TeamAccess returned err: %v", err)
// 	}

// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("TeamAccess is %v, want %v", got, want)
// 	}
// }
