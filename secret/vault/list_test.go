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
	org := "foo"
	repo := "*"
	name := "baz"
	value := "foob"
	typee := "org"
	arr := []string{"foo", "bar"}
	want := []*library.Secret{
		&library.Secret{
			Org:    &org,
			Repo:   &repo,
			Name:   &name,
			Value:  &value,
			Type:   &typee,
			Images: &arr,
			Events: &arr,
		},
	}

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.List(typee, org, repo, 1, 10)

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
	org := "foo"
	repo := "bar"
	name := "baz"
	value := "foob"
	typee := "repo"
	arr := []string{"foo", "bar"}
	want := []*library.Secret{
		&library.Secret{
			Org:    &org,
			Repo:   &repo,
			Name:   &name,
			Value:  &value,
			Type:   &typee,
			Images: &arr,
			Events: &arr,
		},
	}

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
	org := "foo"
	team := "bar"
	name := "baz"
	value := "foob"
	typee := "shared"
	arr := []string{"foo", "bar"}
	want := []*library.Secret{
		&library.Secret{
			Org:    &org,
			Team:   &team,
			Name:   &name,
			Value:  &value,
			Type:   &typee,
			Images: &arr,
			Events: &arr,
		},
	}

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
