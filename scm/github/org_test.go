// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/types/library"
)

func TestGithub_GetOrgName(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/orgs/:org", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/get_org.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	want := "github"

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.GetOrgName(u, "github")

	if resp.Code != http.StatusOK {
		t.Errorf("GetOrgName returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("GetOrgName returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetOrgName is %v, want %v", got, want)
	}
}

func TestGithub_GetOrgName_Personal(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/user.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	want := "octocat"

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.GetOrgName(u, "octocat")

	if resp.Code != http.StatusOK {
		t.Errorf("GetOrgName returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("GetOrgName returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetOrgName is %v, want %v", got, want)
	}
}

func TestGithub_GetOrgName_Fail(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/orgs/:org", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusNotFound)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	u := new(library.User)
	u.SetName("foo")
	u.SetToken("bar")

	client, _ := NewTest(s.URL)

	// run test
	_, err := client.GetOrgName(u, "octocat")

	if err == nil {
		t.Error("GetOrgName should return error")
	}
}
