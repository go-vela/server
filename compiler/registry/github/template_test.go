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

	api "github.com/go-vela/server/api/types"
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
	c, err := New(context.Background(), s.URL, "")
	if err != nil {
		t.Errorf("Creating client returned err: %v", err)
	}

	got, err := c.Template(u, src)

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
	c, err := New(context.Background(), s.URL, "")
	if err != nil {
		t.Errorf("Creating client returned err: %v", err)
	}

	got, err := c.Template(u, src)

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
	c, err := New(context.Background(), s.URL, "")
	if err != nil {
		t.Errorf("Creating client returned err: %v", err)
	}

	got, err := c.Template(u, src)

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
	c, err := New(context.Background(), s.URL, "")
	if err != nil {
		t.Errorf("Creating client returned err: %v", err)
	}

	got, err := c.Template(u, src)

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
	c, err := New(context.Background(), s.URL, "")
	if err != nil {
		t.Errorf("Creating client returned err: %v", err)
	}

	got, err := c.Template(u, src)

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
