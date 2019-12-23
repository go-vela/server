// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
)

func TestVault_Update_Org(t *testing.T) {
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
	engine.PUT("/v1/secret/:type/:org/:name", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	// setup types
	sec := new(library.Secret)
	sec.SetOrg("foo")
	sec.SetRepo("*")
	sec.SetTeam("")
	sec.SetName("bar")
	sec.SetValue("baz")
	sec.SetType("org")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetEvents([]string{"foo", "bar"})

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Update("org", "foo", "*", sec)

	if resp.Code != http.StatusOK {
		t.Errorf("Update returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Update returned err: %v", err)
	}
}

func TestVault_Update_Repo(t *testing.T) {
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
	engine.PUT("/v1/secret/:type/:org/:repo/:name", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	// setup types
	sec := new(library.Secret)
	sec.SetOrg("foo")
	sec.SetRepo("bar")
	sec.SetName("baz")
	sec.SetValue("foob")
	sec.SetType("repo")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetEvents([]string{"foo", "bar"})

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Update("repo", "foo", "bar", sec)

	if resp.Code != http.StatusOK {
		t.Errorf("Update returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Update returned err: %v", err)
	}
}

func TestVault_Update_Shared(t *testing.T) {
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
	engine.PUT("/v1/secret/:type/:org/:team/:name", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	// setup types
	sec := new(library.Secret)
	sec.SetOrg("foo")
	sec.SetTeam("bar")
	sec.SetName("baz")
	sec.SetValue("foob")
	sec.SetType("shared")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetEvents([]string{"foo", "bar"})

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Update("shared", "foo", "bar", sec)

	if resp.Code != http.StatusOK {
		t.Errorf("Update returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Update returned err: %v", err)
	}
}

func TestVault_Update_InvalidSecret(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/v1/secret/:type/:org/:repo/:name", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/invalid_repo.json")
	})
	engine.PUT("/v1/secret/:type/:org/:repo/:name", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	// setup types
	sec := new(library.Secret)
	sec.SetOrg("foo")
	sec.SetRepo("bar")
	sec.SetName("baz")
	sec.SetValue("")
	sec.SetType("repo")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetEvents([]string{"foo", "bar"})

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Update("repo", "foo", "bar", sec)

	if resp.Code != http.StatusOK {
		t.Errorf("Update returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Errorf("Update should have returned err")
	}
}

func TestVault_Update_InvalidType(t *testing.T) {
	// setup types
	sec := new(library.Secret)
	sec.SetOrg("foo")
	sec.SetRepo("bar")
	sec.SetName("baz")
	sec.SetValue("foob")
	sec.SetType("invalid")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetEvents([]string{"foo", "bar"})

	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	defer fake.Close()

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Update("invalid", "foo", "bar", sec)
	if err == nil {
		t.Errorf("Update should have returned err")
	}
}

func TestVault_Update_ClosedServer(t *testing.T) {
	// setup types
	sec := new(library.Secret)
	sec.SetOrg("foo")
	sec.SetRepo("bar")
	sec.SetName("baz")
	sec.SetValue("foob")
	sec.SetType("repo")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetEvents([]string{"foo", "bar"})

	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	fake.Close()

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Update("repo", "foo", "bar", sec)
	if err == nil {
		t.Errorf("Update should have returned err")
	}
}

func TestVault_Update_NoWrite(t *testing.T) {
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
	engine.PUT("/v1/secret/:type/:org/:repo/:name", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	// setup types
	sec := new(library.Secret)
	sec.SetOrg("foo")
	sec.SetRepo("bar")
	sec.SetName("baz")
	sec.SetValue("foob")
	sec.SetType("repo")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetEvents([]string{"foo", "bar"})

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Update("repo", "foo", "bar", sec)

	if resp.Code != http.StatusOK {
		t.Errorf("Update returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Errorf("Update should have returned err")
	}
}
