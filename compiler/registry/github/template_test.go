// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	cacheredis "github.com/go-vela/server/cache/redis"
	"github.com/go-vela/server/compiler/registry"
)

func TestGithub_Template(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:owner/:name/contents/:path", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/template.json")
	})

	s := httptest.NewServer(engine)

	defer s.Close()

	// setup types
	str := "foo"

	r := &api.Repo{
		Org:  new("github"),
		Name: new("octocat"),
	}

	u := &api.User{
		Name:  &str,
		Token: &str,
	}

	src := &registry.Source{
		Org:  "github",
		Repo: "octocat",
		Name: "template.yml",
	}

	want, err := os.ReadFile("testdata/template.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	// run test
	c, err := New(context.Background(), s.URL, "", nil)
	if err != nil {
		t.Errorf("Creating client returned err: %v", err)
	}

	got, err := c.Template(context.Background(), r, u, src, "")

	if resp.Code != http.StatusOK {
		t.Errorf("Template returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Template returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Template is %v, want %v", got, want)
	}
}

func TestGithub_TemplateSourceRef(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// store the ref
	gotRef := ""

	// setup mock server
	engine.GET("/api/v3/repos/:owner/:name/contents/:path", func(c *gin.Context) {
		gotRef = c.Request.FormValue("ref")
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/template.json")
	})

	s := httptest.NewServer(engine)

	defer s.Close()

	// setup types
	str := "foo"

	r := &api.Repo{
		Org:  new("github"),
		Name: new("octocat"),
	}

	u := &api.User{
		Name:  &str,
		Token: &str,
	}

	src := &registry.Source{
		Org:  "github",
		Repo: "octocat",
		Name: "template.yml",
		Ref:  "main",
	}

	want, err := os.ReadFile("testdata/template.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	// run test
	c, err := New(context.Background(), s.URL, "", nil)
	if err != nil {
		t.Errorf("Creating client returned err: %v", err)
	}

	got, err := c.Template(context.Background(), r, u, src, "")

	if resp.Code != http.StatusOK {
		t.Errorf("Template returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Template returned err: %v", err)
	}

	if gotRef != src.Ref {
		t.Errorf("Ref returned %v, want %v", gotRef, src.Ref)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Template is %v, want %v", got, want)
	}
}

func TestGithub_TemplateEmptySourceRef(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// store the ref
	gotRef := ""

	// setup mock server
	engine.GET("/api/v3/repos/:owner/:name/contents/:path", func(c *gin.Context) {
		gotRef = c.Request.FormValue("ref")
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/template.json")
	})

	s := httptest.NewServer(engine)

	defer s.Close()

	// setup types
	str := "foo"

	r := &api.Repo{
		Org:  new("github"),
		Name: new("octocat"),
	}

	u := &api.User{
		Name:  &str,
		Token: &str,
	}

	src := &registry.Source{
		Org:  "github",
		Repo: "octocat",
		Name: "template.yml",
	}

	want, err := os.ReadFile("testdata/template.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	// run test
	c, err := New(context.Background(), s.URL, "", nil)
	if err != nil {
		t.Errorf("Creating client returned err: %v", err)
	}

	got, err := c.Template(context.Background(), r, u, src, "")

	if resp.Code != http.StatusOK {
		t.Errorf("Template returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Template returned err: %v", err)
	}

	if gotRef != "" {
		t.Errorf("Ref returned %v, want empty string", gotRef)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Template is %v, want %v", got, want)
	}
}

func TestGithub_Template_BadRequest(t *testing.T) {
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
	str := "foo"

	r := &api.Repo{
		Org:  new("github"),
		Name: new("octocat"),
	}

	u := &api.User{
		Name:  &str,
		Token: &str,
	}

	src := &registry.Source{
		Org:  "github",
		Repo: "octocat",
		Name: "template.yml",
	}

	// run test
	c, err := New(context.Background(), s.URL, "", nil)
	if err != nil {
		t.Errorf("Creating client returned err: %v", err)
	}

	got, err := c.Template(context.Background(), r, u, src, "")

	if resp.Code != http.StatusOK {
		t.Errorf("Template returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Error("Template should have returned err")
	}

	if got != nil {
		t.Errorf("Template is %v, want nil", got)
	}
}

func TestGithub_Template_NotFound(t *testing.T) {
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
	str := "foo"

	r := &api.Repo{
		Org:  new("github"),
		Name: new("octocat"),
	}

	u := &api.User{
		Name:  &str,
		Token: &str,
	}

	src := &registry.Source{
		Org:  "github",
		Repo: "octocat",
		Name: "template.yml",
	}

	// run test
	c, err := New(context.Background(), s.URL, "", nil)
	if err != nil {
		t.Errorf("Creating client returned err: %v", err)
	}

	got, err := c.Template(context.Background(), r, u, src, "")

	if resp.Code != http.StatusOK {
		t.Errorf("Template returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Error("Template should have returned err")
	}

	if got != nil {
		t.Errorf("Template is %v, want nil", got)
	}
}

func TestGithub_Template_WithCache(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	etag := `"template-etag-123"`
	requestCount := 0

	// setup mock server that supports conditional requests
	engine.GET("/api/v3/repos/:owner/:name/contents/:path", func(c *gin.Context) {
		requestCount++

		if c.GetHeader("If-None-Match") == etag {
			c.Status(http.StatusNotModified)

			return
		}

		c.Header("Content-Type", "application/json")
		c.Header("Etag", etag)
		c.Status(http.StatusOK)
		c.File("testdata/template.json")
	})

	s := httptest.NewServer(engine)

	defer s.Close()

	// setup redis cache
	_cache, err := cacheredis.NewTest("c94bc43c11613ceb6c9f6ac73451e41de90806b2ca6953010b547b20fde9ad90")
	if err != nil {
		t.Errorf("unable to create cache service: %v", err)
	}

	// setup types
	str := "foo"

	r := &api.Repo{
		Org:  new("github"),
		Name: new("octocat"),
	}

	u := &api.User{
		Name:  &str,
		Token: &str,
	}

	src := &registry.Source{
		Org:  "github",
		Repo: "octocat",
		Name: "template.yml",
	}

	want, err := os.ReadFile("testdata/template.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	// run test with cache
	c, err := New(context.Background(), s.URL, "", _cache)
	if err != nil {
		t.Errorf("Creating client returned err: %v", err)
	}

	// first call - cache miss, fetches from server
	got, err := c.Template(context.Background(), r, u, src, "")
	if err != nil {
		t.Errorf("Template (cache miss) returned err: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Template (cache miss) mismatch (-want +got):\n%s", diff)
	}

	// second call - cache hit, server returns 304
	got, err = c.Template(context.Background(), r, u, src, "")
	if err != nil {
		t.Errorf("Template (cache hit) returned err: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Template (cache hit) mismatch (-want +got):\n%s", diff)
	}

	if requestCount != 2 {
		t.Errorf("expected 2 server requests, got %d", requestCount)
	}
}
