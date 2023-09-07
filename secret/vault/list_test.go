// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"context"
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
	engine.GET("/v1/secret/org/foo/:name", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/org.json")
	})
	engine.GET("/v1/secret/org/foo", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/list.json")
	})

	engine.GET("/v1/secret/data/org/foo/:name", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/org.json")
	})
	engine.GET("/v1/secret/metadata/org/foo", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/list.json")
	})

	engine.GET("/v1/secret/data/prefix/org/foo/:name", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/org.json")
	})
	engine.GET("/v1/secret/metadata/prefix/org/foo", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/list.json")
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
			got, err := s.List(context.TODO(), "org", "foo", "*", 1, 10, []string{})

			if resp.Code != http.StatusOK {
				t.Errorf("List returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err != nil {
				t.Errorf("List returned err: %v", err)
			}

			if !reflect.DeepEqual(got, want) {
				t.Errorf("List is %v, want %v", got, want)
			}
		})
	}
}

func TestVault_List_Repo(t *testing.T) {
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
	engine.GET("/v1/secret/repo/foo/bar/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/repo.json")
	})
	engine.GET("/v1/secret/repo/foo/bar/foob", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/repo.json")
	})
	engine.GET("/v1/secret/repo/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/list.json")
	})

	engine.GET("/v1/secret/data/repo/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/repo.json")
	})
	engine.GET("/v1/secret/data/repo/foo/bar/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/repo.json")
	})
	engine.GET("/v1/secret/data/repo/foo/bar/foob", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/repo.json")
	})
	engine.GET("/v1/secret/metadata/repo/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/list.json")
	})

	engine.GET("/v1/secret/data/prefix/repo/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/repo.json")
	})
	engine.GET("/v1/secret/data/prefix/repo/foo/bar/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/repo.json")
	})
	engine.GET("/v1/secret/data/prefix/repo/foo/bar/foob", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/repo.json")
	})
	engine.GET("/v1/secret/metadata/prefix/repo/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/list.json")
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
			got, err := s.List(context.TODO(), "repo", "foo", "bar", 1, 10, []string{})

			if resp.Code != http.StatusOK {
				t.Errorf("List returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err != nil {
				t.Errorf("List returned err: %v", err)
			}

			if !reflect.DeepEqual(got, want) {
				t.Errorf("List is %v, want %v", got, want)
			}
		})
	}
}

func TestVault_List_Shared(t *testing.T) {
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
	engine.GET("/v1/secret/shared/foo/bar/foob", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/shared.json")
	})
	engine.GET("/v1/secret/shared/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/list.json")
	})

	engine.GET("/v1/secret/data/shared/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/shared.json")
	})
	engine.GET("/v1/secret/data/shared/foo/bar/foob", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/shared.json")
	})
	engine.GET("/v1/secret/metadata/shared/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/list.json")
	})

	engine.GET("/v1/secret/data/prefix/shared/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/shared.json")
	})
	engine.GET("/v1/secret/data/prefix/shared/foo/bar/foob", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/shared.json")
	})
	engine.GET("/v1/secret/metadata/prefix/shared/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/list.json")
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
			got, err := s.List(context.TODO(), "shared", "foo", "bar", 1, 10, []string{})

			if resp.Code != http.StatusOK {
				t.Errorf("List returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err != nil {
				t.Errorf("List returned err: %v", err)
			}

			if !reflect.DeepEqual(got, want) {
				t.Errorf("List is %v, want %v", got, want)
			}
		})
	}
}

func TestVault_List_InvalidType(t *testing.T) {
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

			got, err := s.List(context.TODO(), "invalid", "foo", "bar", 1, 10, []string{})
			if err == nil {
				t.Errorf("List should have returned err")
			}

			if got != nil {
				t.Errorf("List is %v, want nil", got)
			}
		})
	}
}

func TestVault_List_ClosedServer(t *testing.T) {
	// setup mock server
	fake := httptest.NewServer(http.NotFoundHandler())
	fake.Close()

	// run test
	s, err := New(
		WithAddress(fake.URL),
		WithAuthMethod(""),
		WithAWSRole(""),
		WithPrefix(""),
		WithToken("foo"),
		WithTokenDuration(0),
		WithVersion("1"),
	)
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	got, err := s.List(context.TODO(), "repo", "foo", "bar", 1, 10, []string{})
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
	engine.GET("/v1/secret/repo/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/empty_list.json")
	})

	engine.GET("/v1/secret/metadata/repo/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/empty_list.json")
	})

	engine.GET("/v1/secret/metadata/prefix/repo/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/empty_list.json")
	})

	fake := httptest.NewServer(engine)
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

			got, err := s.List(context.TODO(), "repo", "foo", "bar", 1, 10, []string{})

			if resp.Code != http.StatusOK {
				t.Errorf("List returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err == nil {
				t.Errorf("List should have returned err")
			}

			if got != nil {
				t.Errorf("List is %v, want nil", got)
			}
		})
	}
}

func TestVault_List_InvalidList(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/v1/secret/repo/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/invalid_list.json")
	})

	engine.GET("/v1/secret/metadata/repo/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/invalid_list.json")
	})

	engine.GET("/v1/secret/metadata/prefix/repo/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/invalid_list.json")
	})

	fake := httptest.NewServer(engine)
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

			got, err := s.List(context.TODO(), "repo", "foo", "bar", 1, 10, []string{})

			if resp.Code != http.StatusOK {
				t.Errorf("List returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err == nil {
				t.Errorf("List should have returned err")
			}

			if got != nil {
				t.Errorf("List is %v, want nil", got)
			}
		})
	}
}

func TestVault_List_NoRead(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/v1/secret/repo/foo/bar/bar", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})
	engine.GET("/v1/secret/repo/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/list.json")
	})

	engine.GET("/v1/secret/data/repo/foo/bar/bar", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})
	engine.GET("/v1/secret/metadata/repo/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/list.json")
	})

	engine.GET("/v1/secret/data/prefix/repo/foo/bar/bar", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})
	engine.GET("/v1/secret/metadata/prefix/repo/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/list.json")
	})

	fake := httptest.NewServer(engine)
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

			got, err := s.List(context.TODO(), "repo", "foo", "bar", 1, 10, []string{})

			if resp.Code != http.StatusOK {
				t.Errorf("List returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err == nil {
				t.Errorf("List should have returned err")
			}

			if got != nil {
				t.Errorf("List is %v, want nil", got)
			}
		})
	}
}
