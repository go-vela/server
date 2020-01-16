// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
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

func TestVault_Create_Org(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
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
	sec.SetAllowCommand(false)

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Create("org", "foo", "*", sec)

	if resp.Code != http.StatusOK {
		t.Errorf("Create returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Create returned err: %v", err)
	}
}

func TestVault_Create_Repo(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.PUT("/v1/secret/:type/:org/:repo/:name", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	// setup types
	sec := new(library.Secret)
	sec.SetOrg("foo")
	sec.SetRepo("bar")
	sec.SetTeam("")
	sec.SetName("baz")
	sec.SetValue("foob")
	sec.SetType("repo")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetEvents([]string{"foo", "bar"})
	sec.SetAllowCommand(false)

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Create("repo", "foo", "bar", sec)

	if resp.Code != http.StatusOK {
		t.Errorf("Create returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Create returned err: %v", err)
	}
}

func TestVault_Create_Shared(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.PUT("/v1/secret/:type/:org/:team/:name", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	// setup types
	sec := new(library.Secret)
	sec.SetOrg("foo")
	sec.SetRepo("")
	sec.SetTeam("bar")
	sec.SetName("baz")
	sec.SetValue("foob")
	sec.SetType("shared")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetEvents([]string{"foo", "bar"})
	sec.SetAllowCommand(false)

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Create("shared", "foo", "bar", sec)

	if resp.Code != http.StatusOK {
		t.Errorf("Create returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Create returned err: %v", err)
	}
}

func TestVault_Create_InvalidSecret(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.PUT("/v1/secret/:type/:org/:repo/:name", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	// setup types
	sec := new(library.Secret)
	sec.SetOrg("foo")
	sec.SetRepo("bar")
	sec.SetTeam("")
	sec.SetName("baz")
	sec.SetValue("")
	sec.SetType("repo")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetEvents([]string{"foo", "bar"})
	sec.SetAllowCommand(false)

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Create("repo", "foo", "bar", sec)

	if resp.Code != http.StatusOK {
		t.Errorf("Create returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Errorf("Create should have returned err")
	}
}

func TestVault_Create_InvalidType(t *testing.T) {
	// setup types
	sec := new(library.Secret)
	sec.SetOrg("foo")
	sec.SetRepo("bar")
	sec.SetTeam("")
	sec.SetName("baz")
	sec.SetValue("foob")
	sec.SetType("invalid")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetEvents([]string{"foo", "bar"})
	sec.SetAllowCommand(false)

	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	defer fake.Close()

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Create("invalid", "foo", "bar", sec)
	if err == nil {
		t.Errorf("Create should have returned err")
	}
}

func TestVault_Create_ClosedServer(t *testing.T) {
	// setup types
	sec := new(library.Secret)
	sec.SetOrg("foo")
	sec.SetRepo("bar")
	sec.SetTeam("")
	sec.SetName("baz")
	sec.SetValue("foob")
	sec.SetType("repo")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetEvents([]string{"foo", "bar"})
	sec.SetAllowCommand(false)

	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	fake.Close()

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Create("repo", "foo", "bar", sec)
	if err == nil {
		t.Errorf("Create should have returned err")
	}
}
