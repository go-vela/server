// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

func TestGithub_Config_YML(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/foo/bar/contents/:path", func(c *gin.Context) {
		if c.Param("path") == ".vela.yaml" {
			c.Status(http.StatusNotFound)
			return
		}

		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/yml.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	want, err := ioutil.ReadFile("testdata/pipeline.yml")
	if err != nil {
		t.Errorf("Config reading file returned err: %v", err)
	}

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Config(u, "foo", "bar", "")

	if resp.Code != http.StatusOK {
		t.Errorf("Config returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Config returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Config is %v, want %v", got, want)
	}
}

func TestGithub_Config_YML_BadRequest(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/foo/bar/contents/:path", func(c *gin.Context) {
		c.Status(http.StatusBadRequest)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Config(u, "foo", "bar", "")

	if resp.Code != http.StatusOK {
		t.Errorf("Config returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Error("Config should have returned err")
	}

	if got != nil {
		t.Errorf("Config is %v, want nil", got)
	}
}

func TestGithub_Config_YAML(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/foo/bar/contents/:path", func(c *gin.Context) {
		if c.Param("path") == ".vela.yml" {
			c.Status(http.StatusNotFound)
			return
		}

		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/yaml.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	want, err := ioutil.ReadFile("testdata/pipeline.yml")
	if err != nil {
		t.Errorf("Config reading file returned err: %v", err)
	}

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Config(u, "foo", "bar", "")

	if resp.Code != http.StatusOK {
		t.Errorf("Config returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Config returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Config is %v, want %v", got, want)
	}
}

func TestGithub_Config_YAML_BadRequest(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/foo/bar/contents/:path", func(c *gin.Context) {
		if c.Param("path") == ".vela.yml" {
			c.Status(http.StatusNotFound)
			return
		}

		c.Status(http.StatusBadRequest)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Config(u, "foo", "bar", "")

	if resp.Code != http.StatusOK {
		t.Errorf("Config returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Error("Config should have returned err")
	}

	if got != nil {
		t.Errorf("Config is %v, want nil", got)
	}
}

func TestGithub_Config_NotFound(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/foo/bar/contents/:path", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Config(u, "foo", "bar", "")

	if resp.Code != http.StatusOK {
		t.Errorf("Config returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Error("Config should have returned err")
	}

	if got != nil {
		t.Errorf("Config is %v, want nil", got)
	}
}

func TestGithub_Disable(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:org/:repo/hooks", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/hooks.json")
	})
	engine.DELETE("/api/v3/repos/:org/:repo/hooks/:hook_id", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL, "https://foo.bar.com")

	// run test
	err := client.Disable(u, "foo", "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Disable returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Disable returned err: %v", err)
	}
}

func TestGithub_Disable_NotFoundHooks(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:org/:repo/hooks", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL, "https://foo.bar.com")

	// run test
	err := client.Disable(u, "foo", "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Disable returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Error("Disable should have returned err")
	}
}

func TestGithub_Disable_HooksButNotFound(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:org/:repo/hooks", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/hooks.json")
	})
	engine.DELETE("/api/v3/repos/:org/:repo/hooks/:hook_id", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL, "https://foos.ball.com")

	// run test
	err := client.Disable(u, "foo", "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Disable returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Disable returned err: %v", err)
	}
}

func TestGithub_Enable(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.POST("/api/v3/repos/:org/:repo/hooks", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/hook.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL)

	// run test
	_, err := client.Enable(u, "foo", "bar", "secret")

	if resp.Code != http.StatusOK {
		t.Errorf("Enable returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Enable returned err: %v", err)
	}
}

func TestGithub_Status_Running(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.POST("/api/v3/repos/:org/:repo/statuses/:sha", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/status.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	b := new(library.Build)
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)
	b.SetEvent(constants.EventPush)
	b.SetStatus(constants.StatusRunning)
	b.SetCommit("abcd1234")

	client, _ := NewTest(s.URL)

	// run test
	err := client.Status(u, b, "foo", "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}
}

func TestGithub_Status_Success(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.POST("/api/v3/repos/:org/:repo/statuses/:sha", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/status.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	b := new(library.Build)
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)
	b.SetEvent(constants.EventPush)
	b.SetStatus(constants.StatusRunning)
	b.SetCommit("abcd1234")

	client, _ := NewTest(s.URL)

	// run test
	err := client.Status(u, b, "foo", "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}
}

func TestGithub_Status_Failure(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.POST("/api/v3/repos/:org/:repo/statuses/:sha", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/status.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	b := new(library.Build)
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)
	b.SetEvent(constants.EventPush)
	b.SetStatus(constants.StatusRunning)
	b.SetCommit("abcd1234")

	client, _ := NewTest(s.URL)

	// run test
	err := client.Status(u, b, "foo", "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}
}

func TestGithub_Status_Killed(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.POST("/api/v3/repos/:org/:repo/statuses/:sha", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/status.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	b := new(library.Build)
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)
	b.SetEvent(constants.EventPush)
	b.SetStatus(constants.StatusRunning)
	b.SetCommit("abcd1234")

	client, _ := NewTest(s.URL)

	// run test
	err := client.Status(u, b, "foo", "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}
}

func TestGithub_Status_Error(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.POST("/api/v3/repos/:org/:repo/statuses/:sha", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/status.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	b := new(library.Build)
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)
	b.SetEvent(constants.EventPush)
	b.SetStatus(constants.StatusRunning)
	b.SetCommit("abcd1234")

	client, _ := NewTest(s.URL)

	// run test
	err := client.Status(u, b, "foo", "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}
}

func TestGithub_ListUserRepos(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/user/repos", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/listuserrepos.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(library.Repo)
	r.SetOrg("octocat")
	r.SetName("Hello-World")
	r.SetFullName("octocat/Hello-World")
	r.SetLink("https://github.com/octocat/Hello-World")
	r.SetClone("https://github.com/octocat/Hello-World.git")
	r.SetBranch("master")
	r.SetPrivate(false)

	want := []*library.Repo{r}

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.ListUserRepos(u)

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repo list is %v, want %v", got, want)
	}
}

func TestGithub_ListUserRepos_Ineligible(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/user/repos", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/listuserrepos_ineligible.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	want := []*library.Repo{}

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.ListUserRepos(u)

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repo list is %v, want %v", got, want)
	}
}

func TestGithub_GetPullRequest(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:owner/:repo/pulls/:pull_number", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/get_pull_request.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(library.Repo)
	r.SetOrg("octocat")
	r.SetName("Hello-World")

	wantCommit := "6dcb09b5b57875f334f61aebed695e2e4193db5e"
	wantBranch := "master"
	wantBaseRef := "master"

	client, _ := NewTest(s.URL)

	// run test
	gotCommit, gotBranch, gotBaseRef, err := client.GetPullRequest(u, r, 1)

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	if !strings.EqualFold(gotCommit, wantCommit) {
		t.Errorf("Commit is %v, want %v", gotCommit, wantCommit)
	}

	if !strings.EqualFold(gotBranch, wantBranch) {
		t.Errorf("Branch is %v, want %v", gotCommit, wantCommit)
	}

	if !strings.EqualFold(gotBaseRef, wantBaseRef) {
		t.Errorf("BaseRef is %v, want %v", gotCommit, wantCommit)
	}
}
