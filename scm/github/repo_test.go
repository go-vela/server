// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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

	want, err := os.ReadFile("testdata/pipeline.yml")
	if err != nil {
		t.Errorf("Config reading file returned err: %v", err)
	}

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(library.Repo)
	r.SetOrg("foo")
	r.SetName("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Config(u, r, "")

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

func TestGithub_ConfigBackoff_YML(t *testing.T) {
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

	want, err := os.ReadFile("testdata/pipeline.yml")
	if err != nil {
		t.Errorf("Config reading file returned err: %v", err)
	}

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(library.Repo)
	r.SetOrg("foo")
	r.SetName("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Config(u, r, "")

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

	r := new(library.Repo)
	r.SetOrg("foo")
	r.SetName("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Config(u, r, "")

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

	want, err := os.ReadFile("testdata/pipeline.yml")
	if err != nil {
		t.Errorf("Config reading file returned err: %v", err)
	}

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(library.Repo)
	r.SetOrg("foo")
	r.SetName("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Config(u, r, "")

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

func TestGithub_Config_Star(t *testing.T) {
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
		c.File("testdata/star.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	want, err := os.ReadFile("testdata/pipeline.yml")
	if err != nil {
		t.Errorf("Config reading file returned err: %v", err)
	}

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(library.Repo)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetPipelineType(constants.PipelineTypeStarlark)

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Config(u, r, "")

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

func TestGithub_Config_Py(t *testing.T) {
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
		c.File("testdata/py.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	want, err := os.ReadFile("testdata/pipeline.yml")
	if err != nil {
		t.Errorf("Config reading file returned err: %v", err)
	}

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(library.Repo)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetPipelineType(constants.PipelineTypeStarlark)

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Config(u, r, "")

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

	r := new(library.Repo)
	r.SetOrg("foo")
	r.SetName("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Config(u, r, "")

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

	r := new(library.Repo)
	r.SetOrg("foo")
	r.SetName("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Config(u, r, "")

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

func TestGithub_Disable_MultipleHooks(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)
	count := 0
	wantCount := 2

	// setup mock server
	engine.GET("/api/v3/repos/:org/:repo/hooks", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/hooks_multi.json")
	})
	engine.DELETE("/api/v3/repos/:org/:repo/hooks/:hook_id", func(c *gin.Context) {
		count++
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

	if count != wantCount {
		t.Errorf("Count returned %d, want %d", count, wantCount)
	}

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

	wantHook := new(library.Hook)
	wantHook.SetWebhookID(1)
	wantHook.SetSourceID("bar-initialize")
	wantHook.SetCreated(1315329987)
	wantHook.SetNumber(1)
	wantHook.SetEvent("initialize")
	wantHook.SetStatus("success")

	r := new(library.Repo)
	r.SetID(1)
	r.SetName("bar")
	r.SetOrg("foo")
	r.SetHash("secret")
	r.SetAllowPush(true)
	r.SetAllowPull(true)
	r.SetAllowDeploy(true)

	client, _ := NewTest(s.URL)

	// run test
	got, _, err := client.Enable(u, r, new(library.Hook))

	if resp.Code != http.StatusOK {
		t.Errorf("Enable returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Enable returned err: %v", err)
	}

	if !reflect.DeepEqual(wantHook, got) {
		t.Errorf("Enable returned hook %v, want %v", got, wantHook)
	}
}

func TestGithub_Update(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.PATCH("/api/v3/repos/:org/:repo/hooks/:hook_id", func(c *gin.Context) {
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

	r := new(library.Repo)
	r.SetID(1)
	r.SetName("bar")
	r.SetOrg("foo")
	r.SetHash("secret")
	r.SetAllowPush(true)
	r.SetAllowPull(true)
	r.SetAllowDeploy(true)

	hookID := int64(1)

	client, _ := NewTest(s.URL)

	// run test
	err := client.Update(u, r, hookID)

	if resp.Code != http.StatusOK {
		t.Errorf("Update returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Update returned err: %v", err)
	}
}

func TestGithub_Status_Deployment(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.POST("/api/v3/repos/:org/:repo/deployments/:deployment/statuses", func(c *gin.Context) {
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
	b.SetEvent(constants.EventDeploy)
	b.SetStatus(constants.StatusRunning)
	b.SetCommit("abcd1234")
	b.SetSource(fmt.Sprintf("%s/%s/%s/deployments/1", s.URL, "foo", "bar"))

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

func TestGithub_Status_Skipped(t *testing.T) {
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
	b.SetStatus(constants.StatusSkipped)
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

func TestGithub_GetRepo(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:owner/:repo", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/get_repo.json")
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

	want := new(library.Repo)
	want.SetOrg("octocat")
	want.SetName("Hello-World")
	want.SetFullName("octocat/Hello-World")
	want.SetLink("https://github.com/octocat/Hello-World")
	want.SetClone("https://github.com/octocat/Hello-World.git")
	want.SetBranch("master")
	want.SetPrivate(false)
	want.SetTopics([]string{"octocat", "atom", "electron", "api"})
	want.SetVisibility("public")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.GetRepo(u, r)

	if resp.Code != http.StatusOK {
		t.Errorf("GetRepo returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("GetRepo returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetRepo is %v, want %v", got, want)
	}
}

func TestGithub_GetRepo_Fail(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:owner/:repo", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusNotFound)
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

	client, _ := NewTest(s.URL)

	// run test
	_, err := client.GetRepo(u, r)

	if err == nil {
		t.Error("GetRepo should return error")
	}
}

func TestGithub_GetOrgAndRepoName(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:owner/:repo", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/get_repo.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	wantOrg := "octocat"
	wantRepo := "Hello-World"

	client, _ := NewTest(s.URL)

	// run test
	gotOrg, gotRepo, err := client.GetOrgAndRepoName(u, "octocat", "Hello-World")

	if resp.Code != http.StatusOK {
		t.Errorf("GetRepoName returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("GetRepoName returned err: %v", err)
	}

	if !reflect.DeepEqual(gotOrg, wantOrg) {
		t.Errorf("GetRepoName org is %v, want %v", gotOrg, wantOrg)
	}

	if !reflect.DeepEqual(gotRepo, wantRepo) {
		t.Errorf("GetRepoName repo is %v, want %v", gotRepo, wantRepo)
	}
}

func TestGithub_GetOrgAndRepoName_Fail(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:owner/:repo", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
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
	_, _, err := client.GetOrgAndRepoName(u, "octocat", "Hello-World")

	if err == nil {
		t.Error("GetRepoName should return error")
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
	r.SetTopics([]string{"octocat", "atom", "electron", "api"})
	r.SetVisibility("public")

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
	wantHeadRef := "new-topic"

	client, _ := NewTest(s.URL)

	// run test
	gotCommit, gotBranch, gotBaseRef, gotHeadRef, err := client.GetPullRequest(u, r, 1)

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	if !strings.EqualFold(gotCommit, wantCommit) {
		t.Errorf("Commit is %v, want %v", gotCommit, wantCommit)
	}

	if !strings.EqualFold(gotBranch, wantBranch) {
		t.Errorf("Branch is %v, want %v", gotBranch, wantBranch)
	}

	if !strings.EqualFold(gotBaseRef, wantBaseRef) {
		t.Errorf("BaseRef is %v, want %v", gotBaseRef, wantBaseRef)
	}

	if !strings.EqualFold(gotHeadRef, wantHeadRef) {
		t.Errorf("HeadRef is %v, want %v", gotHeadRef, wantHeadRef)
	}
}

func TestGithub_GetBranch(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:owner/:repo/branches/:branch", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/branch.json")
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
	r.SetBranch("main")

	wantBranch := "main"
	wantCommit := "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d"

	client, _ := NewTest(s.URL)

	// run test
	gotBranch, gotCommit, err := client.GetBranch(u, r, "main")

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	if !strings.EqualFold(gotBranch, wantBranch) {
		t.Errorf("Branch is %v, want %v", gotBranch, wantBranch)
	}

	if !strings.EqualFold(gotCommit, wantCommit) {
		t.Errorf("Commit is %v, want %v", gotCommit, wantCommit)
	}
}
