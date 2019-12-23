// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
)

func TestVault_List_Org(t *testing.T) {
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
	sec := new(library.Secret)
	sec.SetOrg("foo")
	sec.SetRepo("*")
	sec.SetName("bar")
	sec.SetValue("baz")
	sec.SetType("org")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetEvents([]string{"foo", "bar"})

	want := []*library.Secret{sec}

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.List("org", "foo", "*", 1, 10)

	if resp.Code != http.StatusOK {
		t.Errorf("List returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("List returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("List is %v, want %v", got, want)
	}
}

func TestVault_List_Repo(t *testing.T) {
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
	sec := new(library.Secret)
	sec.SetOrg("foo")
	sec.SetRepo("bar")
	sec.SetName("baz")
	sec.SetValue("foob")
	sec.SetType("repo")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetEvents([]string{"foo", "bar"})

	want := []*library.Secret{sec}

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.List("repo", "foo", "bar", 1, 10)

	if resp.Code != http.StatusOK {
		t.Errorf("List returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("List returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("List is %v, want %v", got, want)
	}
}

func TestVault_List_Shared(t *testing.T) {
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
	sec := new(library.Secret)
	sec.SetOrg("foo")
	sec.SetTeam("bar")
	sec.SetName("baz")
	sec.SetValue("foob")
	sec.SetType("shared")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetEvents([]string{"foo", "bar"})

	want := []*library.Secret{sec}

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.List("shared", "foo", "bar", 1, 10)

	if resp.Code != http.StatusOK {
		t.Errorf("List returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("List returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("List is %v, want %v", got, want)
	}
}

func TestVault_List_InvalidType(t *testing.T) {
	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	defer fake.Close()

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.List("invalid", "foo", "bar", 1, 10)
	if err == nil {
		t.Errorf("List should have returned err")
	}

	if got != nil {
		t.Errorf("List is %v, want nil", got)
	}
}

func TestVault_List_ClosedServer(t *testing.T) {
	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	fake.Close()

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.List("repo", "foo", "bar", 1, 10)
	if err == nil {
		t.Errorf("List should have returned err")
	}

	if got != nil {
		t.Errorf("List is %v, want nil", got)
	}
}

func TestVault_List_EmptyList(t *testing.T) {
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

	got, err := s.List("repo", "foo", "bar", 1, 10)

	if resp.Code != http.StatusOK {
		t.Errorf("List returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Errorf("List should have returned err")
	}

	if got != nil {
		t.Errorf("List is %v, want nil", got)
	}
}

func TestVault_List_InvalidList(t *testing.T) {
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

	got, err := s.List("repo", "foo", "bar", 1, 10)

	if resp.Code != http.StatusOK {
		t.Errorf("List returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Errorf("List should have returned err")
	}

	if got != nil {
		t.Errorf("List is %v, want nil", got)
	}
}

func TestVault_List_NoRead(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/v1/secret/:type/:org/:repo/:name", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})
	engine.GET("/v1/secret/:type/:org/:repo", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/list.json")
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.List("repo", "foo", "bar", 1, 10)

	if resp.Code != http.StatusOK {
		t.Errorf("List returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Errorf("List should have returned err")
	}

	if got != nil {
		t.Errorf("List is %v, want nil", got)
	}
}
