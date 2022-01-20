// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-vela/types/raw"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/types/library"
)

func TestGithub_CreateDeployment(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.POST("/api/v3/repos/:org/:repo/deployments", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/deployment.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(library.Repo)
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	d := new(library.Deployment)
	d.SetID(1)
	d.SetRepoID(1)
	d.SetURL("https://api.github.com/repos/foo/bar/deployments/1")
	d.SetUser("octocat")
	d.SetCommit("a84d88e7554fc1fa21bcbc4efae3c782a70d2b9d")
	d.SetRef("topic-branch")
	d.SetTask("deploy")
	d.SetTarget("production")
	d.SetDescription("Deploy request from Vela")

	client, _ := NewTest(s.URL, "https://foo.bar.com")

	// run test
	err := client.CreateDeployment(u, r, d)

	if resp.Code != http.StatusOK {
		t.Errorf("CreateDeployment returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("CreateDeployment returned err: %v", err)
	}
}

func TestGithub_GetDeployment(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:org/:repo/deployments/:deployment", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/deployment.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

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

	client, _ := NewTest(s.URL, "https://foo.bar.com")

	// run test
	got, err := client.GetDeployment(u, r, 1)

	if resp.Code != http.StatusOK {
		t.Errorf("GetDeployment returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("GetDeployment returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetDeployment is %v, want %v", got, want)
	}
}

func TestGithub_GetDeploymentCount(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:org/:repo/deployments", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/deployments.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(library.Repo)
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	want := int64(2)

	client, _ := NewTest(s.URL, "https://foo.bar.com")

	// run test
	got, err := client.GetDeploymentCount(u, r)

	if resp.Code != http.StatusOK {
		t.Errorf("GetDeployment returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("GetDeployment returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetDeployment is %v, want %v", got, want)
	}
}

func TestGithub_GetDeploymentList(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:org/:repo/deployments", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/deployments.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

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
	d1.SetPayload(nil)

	d2 := new(library.Deployment)
	d2.SetID(2)
	d2.SetRepoID(1)
	d2.SetURL("https://api.github.com/repos/foo/bar/deployments/2")
	d2.SetUser("octocat")
	d2.SetCommit("a84d88e7554fc1fa21bcbc4efae3c782a70d2b9d")
	d2.SetRef("topic-branch")
	d2.SetTask("deploy")
	d2.SetTarget("production")
	d2.SetDescription("Deploy request from Vela")
	d2.SetPayload(raw.StringSliceMap{"deploy": "migrate"})

	want := []*library.Deployment{d2, d1}

	client, _ := NewTest(s.URL, "https://foo.bar.com")

	// run test
	got, err := client.GetDeploymentList(u, r, 1, 100)

	if resp.Code != http.StatusOK {
		t.Errorf("GetDeployment returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("GetDeployment returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetDeployment is %v, want %v", got, want)
	}
}
