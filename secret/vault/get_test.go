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
	org := "foo"
	repo := "*"
	name := "baz"
	value := "foob"
	typee := "org"
	arr := []string{"foo", "bar"}
	want := &library.Secret{
		Org:    &org,
		Repo:   &repo,
		Name:   &name,
		Value:  &value,
		Type:   &typee,
		Images: &arr,
		Events: &arr,
	}

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Get(typee, org, repo, name)

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
	org := "foo"
	repo := "bar"
	name := "baz"
	value := "foob"
	typee := "repo"
	arr := []string{"foo", "bar"}
	want := &library.Secret{
		Org:    &org,
		Repo:   &repo,
		Name:   &name,
		Value:  &value,
		Type:   &typee,
		Images: &arr,
		Events: &arr,
	}

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Get(typee, org, repo, name)

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
	org := "foo"
	team := "bar"
	name := "baz"
	value := "foob"
	typee := "shared"
	arr := []string{"foo", "bar"}
	want := &library.Secret{
		Org:    &org,
		Team:   &team,
		Name:   &name,
		Value:  &value,
		Type:   &typee,
		Images: &arr,
		Events: &arr,
	}

	// run test
	s, err := New(fake.URL, "foo")
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.Get(typee, org, team, name)

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
	s, err := New(fake.URL, "foo")
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
	s, err := New(fake.URL, "foo")
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
