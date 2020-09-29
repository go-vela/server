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
	engine.PUT("/v1/secret/org/foo/bar", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	engine.PUT("/v1/secret/data/org/foo/bar", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	engine.PUT("/v1/secret/data/prefix/org/foo/bar", func(c *gin.Context) {
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
			s, err := New(fake.URL, "foo", tt.args.version, tt.args.prefix,"", "", 0)
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
		})
	}
}

func TestVault_Create_Repo(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.PUT("/v1/secret/repo/foo/bar/baz", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	engine.PUT("/v1/secret/data/repo/foo/bar/baz", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	engine.PUT("/v1/secret/data/prefix/repo/foo/bar/baz", func(c *gin.Context) {
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
			s, err := New(fake.URL, "foo", tt.args.version, tt.args.prefix,"", "", 0)
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
		})
	}
}

func TestVault_Create_Shared(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.PUT("/v1/secret/shared/foo/bar/baz", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})
	engine.PUT("/v1/secret/data/shared/foo/bar/baz", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})
	engine.PUT("/v1/secret/data/prefix/shared/foo/bar/baz", func(c *gin.Context) {
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
			s, err := New(fake.URL, "foo", tt.args.version, tt.args.prefix,"", "", 0)
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
		})
	}
}

func TestVault_Create_InvalidSecret(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.PUT("/v1/secret/repo/foo/bar/baz", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	engine.PUT("/v1/secret/data/repo/foo/bar/baz", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	engine.PUT("/v1/secret/data/prefix/repo/foo/bar/baz", func(c *gin.Context) {
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
			s, err := New(fake.URL, "foo", tt.args.version, tt.args.prefix,"", "", 0)
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
		})
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
			s, err := New(fake.URL, "foo", tt.args.version, tt.args.prefix,"", "", 0)
			if err != nil {
				t.Errorf("New returned err: %v", err)
			}

			err = s.Create("invalid", "foo", "bar", sec)
			if err == nil {
				t.Errorf("Create should have returned err")
			}
		})
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
			s, err := New(fake.URL, "foo", tt.args.version, tt.args.prefix,"", "", 0)
			if err != nil {
				t.Errorf("New returned err: %v", err)
			}

			err = s.Create("repo", "foo", "bar", sec)
			if err == nil {
				t.Errorf("Create should have returned err")
			}
		})
	}
}
