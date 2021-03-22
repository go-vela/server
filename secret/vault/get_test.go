// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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
	engine.GET("/v1/secret/org/foo/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/org.json")
	})

	engine.GET("/v1/secret/data/org/foo/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/org.json")
	})

	engine.GET("/v1/secret/data/prefix/org/foo/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/org.json")
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

	type args struct {
		version string
		prefix  string
	}
	tests := []struct {
		name string
		args args
	}{
		{"v1", args{version: "1", prefix: ""}},
		{"v2", args{version: "2", prefix: ""}},
		{"v2 with prefix", args{version: "2", prefix: "prefix"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New(
				WithAddress(fake.URL),
				WithAuthMethod(""),
				WithAWSRole(""),
				WithPrefix(tt.args.prefix),
				WithToken("foo"),
				WithTokenDuration(0),
				WithVersion(tt.args.version),
			)
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
		})
	}
}

func TestVault_Get_Repo(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/v1/secret/repo/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/repo.json")
	})

	engine.GET("/v1/secret/data/repo/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/repo.json")
	})

	engine.GET("/v1/secret/data/prefix/repo/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/repo.json")
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

	type args struct {
		version string
		prefix  string
	}
	tests := []struct {
		name string
		args args
	}{
		{"v1", args{version: "1", prefix: ""}},
		{"v2", args{version: "2", prefix: ""}},
		{"v2 with prefix", args{version: "2", prefix: "prefix"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New(
				WithAddress(fake.URL),
				WithAuthMethod(""),
				WithAWSRole(""),
				WithPrefix(tt.args.prefix),
				WithToken("foo"),
				WithTokenDuration(0),
				WithVersion(tt.args.version),
			)
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
		})
	}
}

func TestVault_Get_Shared(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/v1/secret/shared/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/shared.json")
	})

	engine.GET("/v1/secret/data/shared/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/shared.json")
	})

	engine.GET("/v1/secret/data/prefix/shared/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/shared.json")
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

	type args struct {
		version string
		prefix  string
	}
	tests := []struct {
		name string
		args args
	}{
		{"v1", args{version: "1", prefix: ""}},
		{"v2", args{version: "2", prefix: ""}},
		{"v2 with prefix", args{version: "2", prefix: "prefix"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New(
				WithAddress(fake.URL),
				WithAuthMethod(""),
				WithAWSRole(""),
				WithPrefix(tt.args.prefix),
				WithToken("foo"),
				WithTokenDuration(0),
				WithVersion(tt.args.version),
			)
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
		})
	}
}

func TestVault_Get_InvalidType(t *testing.T) {
	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	defer fake.Close()

	type args struct {
		version string
		prefix  string
	}
	tests := []struct {
		name string
		args args
	}{
		{"v1", args{version: "1", prefix: ""}},
		{"v2", args{version: "2", prefix: ""}},
		{"v2 with prefix", args{version: "2", prefix: "prefix"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New(
				WithAddress(fake.URL),
				WithAuthMethod(""),
				WithAWSRole(""),
				WithPrefix(tt.args.prefix),
				WithToken("foo"),
				WithTokenDuration(0),
				WithVersion(tt.args.version),
			)
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
		})
	}
}

func TestVault_Get_ClosedServer(t *testing.T) {
	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	fake.Close()

	type args struct {
		version string
		prefix  string
	}
	tests := []struct {
		name string
		args args
	}{
		{"v1", args{version: "1", prefix: ""}},
		{"v2", args{version: "2", prefix: ""}},
		{"v2 with prefix", args{version: "2", prefix: "prefix"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New(
				WithAddress(fake.URL),
				WithAuthMethod(""),
				WithAWSRole(""),
				WithPrefix(tt.args.prefix),
				WithToken("foo"),
				WithTokenDuration(0),
				WithVersion(tt.args.version),
			)
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
		})
	}
}
