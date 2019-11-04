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

func TestVault_Delete_Org(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.DELETE("/v1/secret/:type/:org/:name", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})
	fake := httptest.NewServer(engine)
	defer fake.Close()

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Delete("org", "foo", "bar", "foob")

	if resp.Code != http.StatusOK {
		t.Errorf("Delete returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Delete returned err: %v", err)
	}
}

func TestVault_Delete_Repo(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.DELETE("/v1/secret/:type/:org/:repo/:name", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})
	fake := httptest.NewServer(engine)
	defer fake.Close()

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Delete("repo", "foo", "bar", "foob")

	if resp.Code != http.StatusOK {
		t.Errorf("Delete returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Delete returned err: %v", err)
	}
}

func TestVault_Delete_Shared(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.DELETE("/v1/secret/:type/:org/:team/:name", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})
	fake := httptest.NewServer(engine)
	defer fake.Close()

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Delete("shared", "foo", "bar", "foob")

	if resp.Code != http.StatusOK {
		t.Errorf("Delete returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Delete returned err: %v", err)
	}
}

func TestVault_Delete_InvalidType(t *testing.T) {
	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	defer fake.Close()

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Delete("invalid", "foo", "bar", "foob")
	if err == nil {
		t.Errorf("Delete should have returned err")
	}
}

func TestVault_Delete_ClosedServer(t *testing.T) {
	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	fake.Close()

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	err = s.Delete("repo", "foo", "bar", "foob")
	if err == nil {
		t.Errorf("Delete should have returned err")
	}
}
