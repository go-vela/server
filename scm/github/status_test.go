// SPDX-License-Identifier: Apache-2.0

package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/cache/models"
	"github.com/go-vela/server/constants"
)

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
	r := new(api.Repo)
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")

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
	_, err := client.Status(t.Context(), b, "bar", nil)

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
	r := new(api.Repo)
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")

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
	_, err := client.Status(t.Context(), b, "bar", nil)

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	_, err = client.StepStatus(t.Context(), b, step, "bar", nil)

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}
}

func TestGithub_Status_Running_CheckRun(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.POST("/api/v3/repos/:org/:repo/check-runs", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/check_run.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	r := new(api.Repo)
	r.SetID(1)
	r.SetInstallID(12)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

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

	want := []models.CheckRun{
		{
			ID:          4,
			Context:     "continuous-integration/vela/push",
			Repo:        "foo/bar",
			BuildNumber: 1,
		},
	}

	wantStep := []models.CheckRun{
		{
			ID:          4,
			Context:     "continuous-integration/vela/push/test",
			Repo:        "foo/bar",
			BuildNumber: 1,
		},
	}

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Status(t.Context(), b, "bar", nil)

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Status returned unexpected check runs: %s", diff)
	}

	got, err = client.StepStatus(t.Context(), b, step, "bar", nil)

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	if diff := cmp.Diff(wantStep, got); diff != "" {
		t.Errorf("Status returned unexpected check runs: %s", diff)
	}
}

func TestGithub_Status_PendingApproval_CheckRun_Create(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	var payload map[string]interface{}

	// setup mock server
	engine.POST("/api/v3/repos/:org/:repo/check-runs", func(c *gin.Context) {
		defer c.Request.Body.Close()

		err := json.NewDecoder(c.Request.Body).Decode(&payload)
		if err != nil {
			t.Errorf("unable to decode payload: %v", err)
		}

		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/check_run.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	r := new(api.Repo)
	r.SetID(1)
	r.SetInstallID(12)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)
	b.SetEvent(constants.EventPush)
	b.SetStatus(constants.StatusPendingApproval)
	b.SetCommit("abcd1234")

	client, _ := NewTest(s.URL)

	// run test
	_, err := client.Status(t.Context(), b, "bar", nil)

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	if payload["status"] != "queued" {
		t.Errorf("Status sent %v, want %v", payload["status"], "queued")
	}
}

func TestGithub_Status_Running_CheckRun_Update(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	var payload map[string]interface{}

	// setup mock server
	engine.PATCH("/api/v3/repos/:org/:repo/check-runs/:id", func(c *gin.Context) {
		defer c.Request.Body.Close()

		err := json.NewDecoder(c.Request.Body).Decode(&payload)
		if err != nil {
			t.Errorf("unable to decode payload: %v", err)
		}

		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/check_run.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	r := new(api.Repo)
	r.SetID(1)
	r.SetInstallID(12)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)
	b.SetEvent(constants.EventPush)
	b.SetStatus(constants.StatusRunning)
	b.SetCommit("abcd1234")

	checks := []models.CheckRun{{
		ID:          4,
		Context:     "continuous-integration/vela/push",
		Repo:        "foo/bar",
		BuildNumber: 1,
	}}

	client, _ := NewTest(s.URL)

	// run test
	_, err := client.Status(t.Context(), b, "bar", checks)

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	if payload["status"] != "in_progress" {
		t.Errorf("Status sent %v, want %v", payload["status"], "in_progress")
	}

	if _, ok := payload["conclusion"]; ok {
		t.Errorf("Status sent unexpected conclusion %v", payload["conclusion"])
	}
}

func TestGithub_StepStatus_Success_CheckRun_Update(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	var payload map[string]interface{}

	// setup mock server
	engine.PATCH("/api/v3/repos/:org/:repo/check-runs/:id", func(c *gin.Context) {
		defer c.Request.Body.Close()

		err := json.NewDecoder(c.Request.Body).Decode(&payload)
		if err != nil {
			t.Errorf("unable to decode payload: %v", err)
		}

		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/check_run.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	r := new(api.Repo)
	r.SetID(1)
	r.SetInstallID(12)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)
	b.SetEvent(constants.EventPush)
	b.SetStatus(constants.StatusRunning)
	b.SetCommit("abcd1234")
	b.SetLink("http://localhost/foo/bar/1")

	step := new(api.Step)
	step.SetID(1)
	step.SetNumber(1)
	step.SetName("test")
	step.SetReportAs("test")
	step.SetStatus(constants.StatusSuccess)

	checks := []models.CheckRun{{
		ID:          4,
		Context:     "continuous-integration/vela/push/test",
		Repo:        "foo/bar",
		BuildNumber: 1,
	}}

	client, _ := NewTest(s.URL)

	// run test
	_, err := client.StepStatus(t.Context(), b, step, "bar", checks)

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	if payload["status"] != "completed" {
		t.Errorf("Status sent %v, want %v", payload["status"], "completed")
	}

	if payload["conclusion"] != "success" {
		t.Errorf("Status sent %v, want %v", payload["conclusion"], "success")
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
	r := new(api.Repo)
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")

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
	_, err := client.Status(t.Context(), b, "bar", nil)

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	_, err = client.StepStatus(t.Context(), b, step, "bar", nil)

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}
}

func TestGithub_Status_SuccessMultipleMergeQueue(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	callCount := 0

	// setup mock server
	engine.POST("/api/v3/repos/:org/:repo/statuses/:sha", func(c *gin.Context) {
		callCount++

		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/status.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	r := new(api.Repo)
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetMergeQueueEvents([]string{constants.EventPull, constants.EventPush})

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)
	b.SetEvent(constants.EventMergeGroup)
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
	_, err := client.Status(t.Context(), b, "bar", nil)

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	if callCount != 2 {
		t.Errorf("Expected 2 calls to GitHub API, got %d", callCount)
	}

	callCount = 0

	_, err = client.StepStatus(t.Context(), b, step, "bar", nil)

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	if callCount != 2 {
		t.Errorf("Expected 2 calls to GitHub API, got %d", callCount)
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
	r := new(api.Repo)
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")

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
	_, err := client.Status(t.Context(), b, "bar", nil)

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	_, err = client.StepStatus(t.Context(), b, step, "bar", nil)

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
	r := new(api.Repo)
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")

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
	_, err := client.Status(t.Context(), b, "bar", nil)

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	_, err = client.StepStatus(t.Context(), b, step, "bar", nil)

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
	r := new(api.Repo)
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")

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
	_, err := client.Status(t.Context(), b, "bar", nil)

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	_, err = client.StepStatus(t.Context(), b, step, "bar", nil)

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
	r := new(api.Repo)
	r.SetID(1)
	r.SetOrg("foo")
	r.SetName("bar")

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
	_, err := client.Status(t.Context(), b, "bar", nil)

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}

	_, err = client.StepStatus(t.Context(), b, step, "bar", nil)

	if resp.Code != http.StatusOK {
		t.Errorf("Status returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Status returned err: %v", err)
	}
}
