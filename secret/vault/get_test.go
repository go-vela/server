// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
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

func TestVault_Get_Org(t *testing.T) {
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

	fake := httptest.NewServer(engine)
	defer fake.Close()

	// setup types
	want := new(library.Secret)
	want.SetOrg("foo")
	want.SetRepo("*")
	want.SetName("bar")
	want.SetValue("baz")
	want.SetType("org")
	want.SetImages([]string{"foo", "bar"})
	want.SetEvents([]string{"foo", "bar"})

	// run test
	s, err := New(fake.URL, "foo", "1")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Get("org", "foo", "bar", "baz")

	if resp.Code != http.StatusOK {
		t.Errorf("Get returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Get returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Get is %v, want %v", got, want)
	}
}

func TestVault_Get_Repo(t *testing.T) {
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

	fake := httptest.NewServer(engine)
	defer fake.Close()

	// setup types
	want := new(library.Secret)
	want.SetOrg("foo")
	want.SetRepo("bar")
	want.SetName("baz")
	want.SetValue("foob")
	want.SetType("repo")
	want.SetImages([]string{"foo", "bar"})
	want.SetEvents([]string{"foo", "bar"})

	// run test
	s, err := New(fake.URL, "foo", "1")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Get("repo", "foo", "bar", "baz")

	if resp.Code != http.StatusOK {
		t.Errorf("Get returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Get returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Get is %v, want %v", got, want)
	}
}

func TestVault_Get_Shared(t *testing.T) {
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

	fake := httptest.NewServer(engine)
	defer fake.Close()

	// setup types
	want := new(library.Secret)
	want.SetOrg("foo")
	want.SetTeam("bar")
	want.SetName("baz")
	want.SetValue("foob")
	want.SetType("shared")
	want.SetImages([]string{"foo", "bar"})
	want.SetEvents([]string{"foo", "bar"})

	// run test
	s, err := New(fake.URL, "foo", "1")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Get("shared", "foo", "bar", "baz")

	if resp.Code != http.StatusOK {
		t.Errorf("Get returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Get returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Get is %v, want %v", got, want)
	}
}

func TestVault_Get_InvalidType(t *testing.T) {
	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	defer fake.Close()

	// run test
	s, err := New(fake.URL, "foo", "1")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Get("invalid", "foo", "bar", "foob")
	if err == nil {
		t.Errorf("Get should have returned err")
	}

	if got != nil {
		t.Errorf("Get is %v, want nil", got)
	}
}

func TestVault_Get_ClosedServer(t *testing.T) {
	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	fake.Close()

	// run test
	s, err := New(fake.URL, "foo", "1")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Get("repo", "foo", "bar", "foob")
	if err == nil {
		t.Errorf("Get should have returned err")
	}

	if got != nil {
		t.Errorf("Get is %v, want nil", got)
	}
}
