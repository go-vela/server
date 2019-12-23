// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestVault_Count_Org(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/v1/secret/:type/:org/:name", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/org.json")
	})
	engine.GET("/v1/secret/:type/:org", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/list.json")
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	// setup types
	want := 1

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Count("org", "foo", "*")

	if resp.Code != http.StatusOK {
		t.Errorf("Count returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Count returned err: %v", err)
	}

	if got != int64(want) {
		t.Errorf("Count is %v, want %v", got, want)
	}
}

func TestVault_Count_Repo(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/v1/secret/:type/:org/:repo/:name", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/repo.json")
	})
	engine.GET("/v1/secret/:type/:org/:repo", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/list.json")
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	// setup types
	want := 1

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Count("repo", "foo", "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Count returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Count returned err: %v", err)
	}

	if got != int64(want) {
		t.Errorf("Count is %v, want %v", got, want)
	}
}

func TestVault_Count_Shared(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/v1/secret/:type/:org/:team/:name", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/shared.json")
	})
	engine.GET("/v1/secret/:type/:org/:team", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/list.json")
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	// setup types
	want := 1

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Count("shared", "foo", "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("List returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("List returned err: %v", err)
	}

	if got != int64(want) {
		t.Errorf("Count is %v, want %v", got, want)
	}
}

func TestVault_Count_InvalidType(t *testing.T) {
	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	defer fake.Close()

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Count("invalid", "foo", "bar")
	if err == nil {
		t.Errorf("Count should have returned err")
	}

	if got != 0 {
		t.Errorf("Count is %v, want 0", got)
	}
}

func TestVault_Count_ClosedServer(t *testing.T) {
	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	fake.Close()

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Count("repo", "foo", "bar")
	if err == nil {
		t.Errorf("Count should have returned err")
	}

	if got != 0 {
		t.Errorf("Count is %v, want 0", got)
	}
}

func TestVault_Count_EmptyList(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/v1/secret/:type/:org/:team", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/empty_list.json")
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Count("repo", "foo", "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Count returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Errorf("Count should have returned err")
	}

	if got != 0 {
		t.Errorf("Count is %v, want 0", got)
	}
}

func TestVault_Count_InvalidList(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/v1/secret/:type/:org/:team", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/invalid_list.json")
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Count("repo", "foo", "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Count returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Errorf("Count should have returned err")
	}

	if got != 0 {
		t.Errorf("Count is %v, want 0", got)
	}
}
