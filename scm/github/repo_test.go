// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v74/github"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/yaml/yaml"
	"github.com/go-vela/server/constants"
)

func TestGithub_Config_YML(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/foo/bar/contents/:path", func(c *gin.Context) {
		if c.Param("path") == ".vela.yml" {
			c.Header("Content-Type", "application/json")
			c.Status(http.StatusOK)
			c.File("testdata/yml.json")
			return
		}

		c.Status(http.StatusNotFound)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	want, err := os.ReadFile("testdata/pipeline.yml")
	if err != nil {
		t.Errorf("Config reading file returned err: %v", err)
	}

	// setup types
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetOrg("foo")
	r.SetName("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Config(context.TODO(), u, r, "")

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

	// counter for api calls
	count := 0

	// setup mock server
	engine.GET("/api/v3/repos/foo/bar/contents/:path", func(c *gin.Context) {
		// load the yml file on the second api call
		if c.Param("path") == ".vela.yml" && count != 0 {
			c.Header("Content-Type", "application/json")
			c.Status(http.StatusOK)
			c.File("testdata/yml.json")
			return
		}

		c.Status(http.StatusNotFound)

		// increment api call counter
		count++
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	want, err := os.ReadFile("testdata/pipeline.yml")
	if err != nil {
		t.Errorf("ConfigBackoff reading file returned err: %v", err)
	}

	// setup types
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetOrg("foo")
	r.SetName("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.ConfigBackoff(context.TODO(), u, r, "")

	if resp.Code != http.StatusOK {
		t.Errorf("ConfigBackoff returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("ConfigBackoff returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ConfigBackoff is %v, want %v", got, want)
	}
}

func TestGithub_Config_YAML(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/foo/bar/contents/:path", func(c *gin.Context) {
		if c.Param("path") == ".vela.yaml" {
			c.Header("Content-Type", "application/json")
			c.Status(http.StatusOK)
			c.File("testdata/yaml.json")
			return
		}

		c.Status(http.StatusNotFound)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	want, err := os.ReadFile("testdata/pipeline.yml")
	if err != nil {
		t.Errorf("Config reading file returned err: %v", err)
	}

	// setup types
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetOrg("foo")
	r.SetName("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Config(context.TODO(), u, r, "")

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
		if c.Param("path") == ".vela.star" {
			c.Header("Content-Type", "application/json")
			c.Status(http.StatusOK)
			c.File("testdata/star.json")
			return
		}

		c.Status(http.StatusNotFound)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	want, err := os.ReadFile("testdata/pipeline.star")
	if err != nil {
		t.Errorf("Config reading file returned err: %v", err)
	}

	// setup types
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetPipelineType(constants.PipelineTypeStarlark)

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Config(context.TODO(), u, r, "")

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

func TestGithub_Config_Star_Prefer(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/foo/bar/contents/:path", func(c *gin.Context) {
		// repo has .vela.yml and .vela.star
		switch c.Param("path") {
		case ".vela.yml":
			c.Header("Content-Type", "application/json")
			c.Status(http.StatusOK)
			c.File("testdata/yml.json")
		case ".vela.star":
			c.Header("Content-Type", "application/json")
			c.Status(http.StatusOK)
			c.File("testdata/star.json")
		default:
			c.Status(http.StatusNotFound)
		}
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	want, err := os.ReadFile("testdata/pipeline.star")
	if err != nil {
		t.Errorf("Config reading file returned err: %v", err)
	}

	// setup types
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetPipelineType(constants.PipelineTypeStarlark)

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Config(context.TODO(), u, r, "")

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
		if c.Param("path") == ".vela.py" {
			c.Header("Content-Type", "application/json")
			c.Status(http.StatusOK)
			c.File("testdata/py.json")
			return
		}

		c.Status(http.StatusNotFound)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	want, err := os.ReadFile("testdata/pipeline.star")
	if err != nil {
		t.Errorf("Config reading file returned err: %v", err)
	}

	// setup types
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetPipelineType(constants.PipelineTypeStarlark)

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Config(context.TODO(), u, r, "")

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
		// first default not found
		if c.Param("path") == ".vela.yml" {
			c.Status(http.StatusNotFound)
			return
		}

		// second default (.vela.yaml) causes bad request
		c.Status(http.StatusBadRequest)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetOrg("foo")
	r.SetName("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Config(context.TODO(), u, r, "")

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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetOrg("foo")
	r.SetName("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Config(context.TODO(), u, r, "")

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

func TestGithub_Config_BadEncoding(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/foo/bar/contents/:path", func(c *gin.Context) {
		if c.Param("path") == ".vela.yml" {
			c.Header("Content-Type", "application/json")
			c.Status(http.StatusOK)
			c.File("testdata/yml_bad_encoding.json")
			return
		}

		c.Status(http.StatusNotFound)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetOrg("foo")
	r.SetName("bar")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Config(context.TODO(), u, r, "")

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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL, "https://foo.bar.com")

	// run test
	err := client.Disable(context.TODO(), u, "foo", "bar")

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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL, "https://foo.bar.com")

	// run test
	err := client.Disable(context.TODO(), u, "foo", "bar")

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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL, "https://foos.ball.com")

	// run test
	err := client.Disable(context.TODO(), u, "foo", "bar")

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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL, "https://foo.bar.com")

	// run test
	err := client.Disable(context.TODO(), u, "foo", "bar")

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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	wantHook := new(api.Hook)
	wantHook.SetWebhookID(1)
	wantHook.SetSourceID("bar-initialize")
	wantHook.SetCreated(1315329987)
	wantHook.SetNumber(1)
	wantHook.SetEvent("initialize")
	wantHook.SetStatus("success")

	r := new(api.Repo)
	r.SetID(1)
	r.SetName("bar")
	r.SetOrg("foo")
	r.SetHash("secret")

	client, _ := NewTest(s.URL)

	// run test
	got, _, err := client.Enable(context.TODO(), u, r, new(api.Hook))

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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetID(1)
	r.SetName("bar")
	r.SetOrg("foo")
	r.SetHash("secret")

	hookID := int64(1)

	client, _ := NewTest(s.URL)

	// run test
	_, err := client.Update(context.TODO(), u, r, hookID)

	if resp.Code != http.StatusOK {
		t.Errorf("Update returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Update returned err: %v", err)
	}
}

func TestGithub_Update_webhookExists_True(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.PATCH("/api/v3/repos/:org/:repo/hooks/:hook_id", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)

	client, _ := NewTest(s.URL)

	// run test
	webhookExists, err := client.Update(context.TODO(), u, r, 0)

	if !webhookExists {
		t.Errorf("Update returned %v, want %v", webhookExists, true)
	}

	if err != nil {
		t.Errorf("Update returned err: %v", err)
	}
}

func TestGithub_Update_webhookExists_False(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.PATCH("/api/v3/repos/:org/:repo/hooks/:hook_id", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusNotFound)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)

	client, _ := NewTest(s.URL)

	// run test
	webhookExists, err := client.Update(context.TODO(), u, r, 0)

	if webhookExists {
		t.Errorf("Update returned %v, want %v", webhookExists, false)
	}

	if err == nil {
		t.Error("Update should return error")
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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetID(1)

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)
	b.SetEvent(constants.EventDeploy)
	b.SetStatus(constants.StatusRunning)
	b.SetCommit("abcd1234")
	b.SetSource(fmt.Sprintf("%s/%s/%s/deployments/1", s.URL, "foo", "bar"))

	client, _ := NewTest(s.URL)

	// run test
	err := client.Status(context.TODO(), u, b, "foo", "bar")

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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetID(1)

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)
	b.SetEvent(constants.EventPush)
	b.SetStatus(constants.StatusRunning)
	b.SetCommit("abcd1234")

	step := new(api.Step)
	step.SetID(1)
	step.SetNumber(1)
	step.SetName("test")
	step.SetReportAs("test")
	step.SetStatus(constants.StatusRunning)

	client, _ := NewTest(s.URL)

	// run test
	err := client.Status(context.TODO(), u, b, "foo", "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	err = client.StepStatus(context.TODO(), u, b, step, "foo", "bar")

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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetID(1)

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)
	b.SetEvent(constants.EventPush)
	b.SetStatus(constants.StatusRunning)
	b.SetCommit("abcd1234")

	step := new(api.Step)
	step.SetID(1)
	step.SetNumber(1)
	step.SetName("test")
	step.SetReportAs("test")
	step.SetStatus(constants.StatusSuccess)

	client, _ := NewTest(s.URL)

	// run test
	err := client.Status(context.TODO(), u, b, "foo", "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	err = client.StepStatus(context.TODO(), u, b, step, "foo", "bar")

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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetID(1)

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)
	b.SetEvent(constants.EventPush)
	b.SetStatus(constants.StatusRunning)
	b.SetCommit("abcd1234")

	step := new(api.Step)
	step.SetID(1)
	step.SetNumber(1)
	step.SetName("test")
	step.SetReportAs("test")
	step.SetStatus(constants.StatusFailure)

	client, _ := NewTest(s.URL)

	// run test
	err := client.Status(context.TODO(), u, b, "foo", "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	err = client.StepStatus(context.TODO(), u, b, step, "foo", "bar")

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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetID(1)

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)
	b.SetEvent(constants.EventPush)
	b.SetStatus(constants.StatusRunning)
	b.SetCommit("abcd1234")

	step := new(api.Step)
	step.SetID(1)
	step.SetNumber(1)
	step.SetName("test")
	step.SetReportAs("test")
	step.SetStatus(constants.StatusKilled)

	client, _ := NewTest(s.URL)

	// run test
	err := client.Status(context.TODO(), u, b, "foo", "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	err = client.StepStatus(context.TODO(), u, b, step, "foo", "bar")

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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetID(1)

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)
	b.SetEvent(constants.EventPush)
	b.SetStatus(constants.StatusSkipped)
	b.SetCommit("abcd1234")

	step := new(api.Step)
	step.SetID(1)
	step.SetNumber(1)
	step.SetName("test")
	step.SetReportAs("test")
	step.SetStatus(constants.StatusSkipped)

	client, _ := NewTest(s.URL)

	// run test
	err := client.Status(context.TODO(), u, b, "foo", "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	err = client.StepStatus(context.TODO(), u, b, step, "foo", "bar")

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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetID(1)

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)
	b.SetEvent(constants.EventPush)
	b.SetStatus(constants.StatusRunning)
	b.SetCommit("abcd1234")

	step := new(api.Step)
	step.SetID(1)
	step.SetNumber(1)
	step.SetName("test")
	step.SetReportAs("test")
	step.SetStatus(constants.StatusError)

	client, _ := NewTest(s.URL)

	// run test
	err := client.Status(context.TODO(), u, b, "foo", "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	err = client.StepStatus(context.TODO(), u, b, step, "foo", "bar")

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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetOrg("octocat")
	r.SetName("Hello-World")

	want := new(api.Repo)
	want.SetOrg("octocat")
	want.SetName("Hello-World")
	want.SetFullName("octocat/Hello-World")
	want.SetLink("https://github.com/octocat/Hello-World")
	want.SetClone("https://github.com/octocat/Hello-World.git")
	want.SetBranch("main")
	want.SetPrivate(false)
	want.SetTopics([]string{"octocat", "atom", "electron", "api"})
	want.SetVisibility("public")
	want.SetCustomProps(map[string]any{
		"prop_1": "foo",
		"prop_2": "bar",
	})

	client, _ := NewTest(s.URL)

	// run test
	got, code, err := client.GetRepo(context.TODO(), u, r)

	if code != http.StatusOK {
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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetOrg("octocat")
	r.SetName("Hello-World")

	client, _ := NewTest(s.URL)

	// run test
	_, code, err := client.GetRepo(context.TODO(), u, r)

	if err == nil {
		t.Error("GetRepo should return error")
	}

	if code != http.StatusNotFound {
		t.Errorf("GetRepo should have returned %d status, got %d", http.StatusNotFound, code)
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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	wantOrg := "octocat"
	wantRepo := "Hello-World"

	client, _ := NewTest(s.URL)

	// run test
	gotOrg, gotRepo, err := client.GetOrgAndRepoName(context.TODO(), u, "octocat", "Hello-World")

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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL)

	// run test
	_, _, err := client.GetOrgAndRepoName(context.TODO(), u, "octocat", "Hello-World")

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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetOrg("octocat")
	r.SetName("Hello-World")
	r.SetFullName("octocat/Hello-World")
	r.SetLink("https://github.com/octocat/Hello-World")
	r.SetClone("https://github.com/octocat/Hello-World.git")
	r.SetBranch("main")
	r.SetPrivate(false)
	r.SetTopics([]string{"octocat", "atom", "electron", "api"})
	r.SetVisibility("public")

	want := []string{"octocat/Hello-World"}

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.ListUserRepos(context.TODO(), u)
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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	want := []string{}

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.ListUserRepos(context.TODO(), u)
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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetOrg("octocat")
	r.SetName("Hello-World")
	r.SetOwner(u)

	wantCommit := "6dcb09b5b57875f334f61aebed695e2e4193db5e"
	wantBranch := "main"
	wantBaseRef := "main"
	wantHeadRef := "new-topic"

	client, _ := NewTest(s.URL)

	// run test
	gotCommit, gotBranch, gotBaseRef, gotHeadRef, err := client.GetPullRequest(context.TODO(), r, 1)
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
	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	r := new(api.Repo)
	r.SetOrg("octocat")
	r.SetName("Hello-World")
	r.SetFullName("octocat/Hello-World")
	r.SetBranch("main")
	r.SetOwner(u)

	wantBranch := "main"
	wantCommit := "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d"

	client, _ := NewTest(s.URL)

	// run test
	gotBranch, gotCommit, err := client.GetBranch(context.TODO(), r, "main")
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

func TestGithub_GetNetrcPassword(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/app/installations", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/installations.json")
	})
	engine.POST("/api/v3/app/installations/:id/access_tokens", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/installations_access_tokens.json")
	})
	engine.GET("/api/v3/repos/:org/:repo/collaborators/foo/permission", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/repo_admin.json")
	})
	engine.GET("/api/v3/repos/:org/:repo/collaborators/charlatan/permission", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusForbidden)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	installedRepo := new(api.Repo)
	installedRepo.SetOrg("octocat")
	installedRepo.SetName("Hello-World")
	installedRepo.SetInstallID(1)

	otherRepo := new(api.Repo)
	otherRepo.SetOrg("octocat")
	otherRepo.SetName("Hi-World")
	otherRepo.SetInstallID(2)

	oauthRepo := new(api.Repo)
	oauthRepo.SetOrg("octocat")
	oauthRepo.SetName("Hello-World2")
	oauthRepo.SetInstallID(0)

	u := new(api.User)
	u.SetName("foo")
	u.SetToken("bar")

	badUser := new(api.User)
	badUser.SetName("charlatan")
	badUser.SetToken("bar")

	tests := []struct {
		name          string
		repo          *api.Repo
		user          *api.User
		git           yaml.Git
		appsTransport bool
		wantToken     string
		wantErr       bool
	}{
		{
			name: "installation token",
			repo: installedRepo,
			user: u,
			git: yaml.Git{
				Token: yaml.Token{
					Repositories: []string{"Hello-World"},
					Permissions:  map[string]string{"contents": "read"},
				},
			},
			appsTransport: true,
			wantToken:     "ghs_16C7e42F292c6912E7710c838347Ae178B4a",
			wantErr:       false,
		},
		{
			name: "no app configured returns user oauth token",
			repo: installedRepo,
			user: u,
			git: yaml.Git{
				Token: yaml.Token{
					Repositories: []string{"Hello-World"},
					Permissions:  map[string]string{"contents": "read"},
				},
			},
			appsTransport: false,
			wantToken:     "bar",
			wantErr:       false,
		},
		{
			name: "repo not installed returns user oauth token",
			repo: oauthRepo,
			user: u,
			git: yaml.Git{
				Token: yaml.Token{
					Repositories: []string{"Hello-World"},
					Permissions:  map[string]string{"contents": "read"},
				},
			},
			appsTransport: true,
			wantToken:     "bar",
			wantErr:       false,
		},
		{
			name: "invalid permission resource",
			repo: installedRepo,
			user: u,
			git: yaml.Git{
				Token: yaml.Token{
					Repositories: []string{"Hello-World"},
					Permissions:  map[string]string{"invalid": "read"},
				},
			},
			appsTransport: true,
			wantToken:     "bar",
			wantErr:       true,
		},
		{
			name: "invalid permission level",
			repo: installedRepo,
			user: u,
			git: yaml.Git{
				Token: yaml.Token{
					Repositories: []string{"Hello-World"},
					Permissions:  map[string]string{"contents": "invalid"},
				},
			},
			appsTransport: true,
			wantToken:     "bar",
			wantErr:       true,
		},
		{
			name: "owner with inadequate permission to other repo",
			repo: otherRepo,
			user: badUser,
			git: yaml.Git{
				Token: yaml.Token{
					Repositories: []string{"Hello-World"},
				},
			},
			appsTransport: true,
			wantToken:     "bar",
			wantErr:       true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, _ := NewTest(s.URL)
			if test.appsTransport {
				client.AppsTransport = NewTestAppsTransport(s.URL)
			}

			got, err := client.GetNetrcPassword(context.TODO(), nil, test.repo, test.user, test.git)
			if (err != nil) != test.wantErr {
				t.Errorf("GetNetrcPassword() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.wantToken {
				t.Errorf("GetNetrcPassword() = %v, want %v", got, test.wantToken)
			}
		})
	}
}

func TestGithub_SyncRepoWithInstallation(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/app/installations", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/installations.json")
	})
	engine.POST("/api/v3/app/installations/:id/access_tokens", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/installations_access_tokens.json")
	})
	engine.GET("/api/v3/installation/repositories", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/installation_repositories.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	tests := []struct {
		name           string
		org            string
		repo           string
		wantInstallID  int64
		wantStatusCode int
	}{
		{
			name:           "match",
			org:            "octocat",
			repo:           "Hello-World",
			wantInstallID:  1,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "no match",
			repo:           "octocat/Hello-World2",
			wantInstallID:  0,
			wantStatusCode: http.StatusOK,
		},
	}
	for _, test := range tests {
		// setup types
		r := new(api.Repo)
		r.SetOrg(test.org)
		r.SetName(test.repo)
		r.SetFullName(fmt.Sprintf("%s/%s", test.org, test.repo))

		client, _ := NewTest(s.URL)
		client.AppsTransport = NewTestAppsTransport(s.URL)

		// run test
		got, err := client.SyncRepoWithInstallation(context.TODO(), r)

		if resp.Code != test.wantStatusCode {
			t.Errorf("SyncRepoWithInstallation %s returned %v, want %v", test.name, resp.Code, http.StatusOK)
		}

		if err != nil {
			t.Errorf("SyncRepoWithInstallation %s returned err: %v", test.name, err)
		}

		if got.GetInstallID() != test.wantInstallID {
			t.Errorf("SyncRepoWithInstallation %s returned %v, want %v", test.name, got.GetInstallID(), test.wantInstallID)
		}
	}
}

func TestGithub_applyGitHubInstallationPermission(t *testing.T) {
	tests := []struct {
		name      string
		perms     *github.InstallationPermissions
		resource  string
		perm      string
		wantPerms *github.InstallationPermissions
		wantErr   bool
	}{
		{
			name: "valid read permission for contents",
			perms: &github.InstallationPermissions{
				Contents: github.Ptr(AppInstallPermissionNone),
			},
			resource: AppInstallResourceContents,
			perm:     AppInstallPermissionRead,
			wantPerms: &github.InstallationPermissions{
				Contents: github.Ptr(AppInstallPermissionRead),
			},
			wantErr: false,
		},
		{
			name: "valid write permission for checks",
			perms: &github.InstallationPermissions{
				Checks: github.Ptr(AppInstallPermissionNone),
			},
			resource: AppInstallResourceChecks,
			perm:     AppInstallPermissionWrite,
			wantPerms: &github.InstallationPermissions{
				Checks: github.Ptr(AppInstallPermissionWrite),
			},
			wantErr: false,
		},
		{
			name: "invalid permission value",
			perms: &github.InstallationPermissions{
				Contents: github.Ptr(AppInstallPermissionNone),
			},
			resource: AppInstallResourceContents,
			perm:     "invalid",
			wantPerms: &github.InstallationPermissions{
				Contents: github.Ptr(AppInstallPermissionNone),
			},
			wantErr: true,
		},
		{
			name: "invalid permission key",
			perms: &github.InstallationPermissions{
				Contents: github.Ptr(AppInstallPermissionNone),
			},
			resource: "invalid",
			perm:     AppInstallPermissionRead,
			wantPerms: &github.InstallationPermissions{
				Contents: github.Ptr(AppInstallPermissionNone),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ApplyInstallationPermissions(tt.resource, tt.perm, tt.perms)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToGitHubAppInstallationPermissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.wantPerms, got); diff != "" {
				t.Errorf("ToGitHubAppInstallationPermissions() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
