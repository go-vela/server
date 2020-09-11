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

func TestVault_Update_Org(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/v1/secret/org/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/org.json")
	})
	engine.PUT("/v1/secret/org/foo/bar", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	engine.GET("/v1/secret/data/org/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/org.json")
	})
	engine.PUT("/v1/secret/data/org/foo/bar", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	engine.GET("/v1/secret/data/prefix/org/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/org.json")
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
			s, err := New(fake.URL, "foo", tt.args.version, tt.args.prefix)
			if err != nil {
				t.Errorf("New returned err: %v", err)
			}

			err = s.Update("org", "foo", "*", sec)

			if resp.Code != http.StatusOK {
				t.Errorf("Update returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err != nil {
				t.Errorf("Update returned err: %v", err)
			}
		})
	}
}

func TestVault_Update_Repo(t *testing.T) {
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
	engine.PUT("/v1/secret/repo/foo/bar/baz", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	engine.GET("/v1/secret/data/repo/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/repo.json")
	})
	engine.PUT("/v1/secret/data/repo/foo/bar/baz", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	engine.GET("/v1/secret/data/prefix/repo/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/repo.json")
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
			s, err := New(fake.URL, "foo", tt.args.version, tt.args.prefix)
			if err != nil {
				t.Errorf("New returned err: %v", err)
			}

			err = s.Update("repo", "foo", "bar", sec)

			if resp.Code != http.StatusOK {
				t.Errorf("Update returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err != nil {
				t.Errorf("Update returned err: %v", err)
			}
		})
	}
}

func TestVault_Update_Shared(t *testing.T) {
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
	engine.PUT("/v1/secret/shared/foo/bar/baz", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	engine.GET("/v1/secret/data/shared/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/shared.json")
	})
	engine.PUT("/v1/secret/data/shared/foo/bar/baz", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	engine.GET("/v1/secret/data/prefix/shared/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/shared.json")
	})
	engine.PUT("/v1/secret/data/prefix/shared/foo/bar/baz", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
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
			s, err := New(fake.URL, "foo", tt.args.version, tt.args.prefix)
			if err != nil {
				t.Errorf("New returned err: %v", err)
			}

			err = s.Update("shared", "foo", "bar", sec)

			if resp.Code != http.StatusOK {
				t.Errorf("Update returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err != nil {
				t.Errorf("Update returned err: %v", err)
			}
		})
	}
}

func TestVault_Update_InvalidSecret(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/v1/secret/repo/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/invalid_repo.json")
	})
	engine.PUT("/v1/secret/repo/foo/bar/baz", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	engine.GET("/v1/secret/data/repo/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/invalid_repo.json")
	})
	engine.PUT("/v1/secret/data/repo/foo/bar/baz", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	engine.GET("/v1/secret/data/prefix/repo/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/invalid_repo.json")
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
			s, err := New(fake.URL, "foo", tt.args.version, tt.args.prefix)
			if err != nil {
				t.Errorf("New returned err: %v", err)
			}

			err = s.Update("repo", "foo", "bar", sec)

			if resp.Code != http.StatusOK {
				t.Errorf("Update returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err == nil {
				t.Errorf("Update should have returned err")
			}
		})
	}
}

func TestVault_Update_InvalidType(t *testing.T) {
	// setup types
	sec := new(library.Secret)
	sec.SetOrg("foo")
	sec.SetRepo("bar")
	sec.SetName("baz")
	sec.SetValue("foob")
	sec.SetType("invalid")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetEvents([]string{"foo", "bar"})

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
			s, err := New(fake.URL, "foo", tt.args.version, tt.args.prefix)
			if err != nil {
				t.Errorf("New returned err: %v", err)
			}

			err = s.Update("invalid", "foo", "bar", sec)
			if err == nil {
				t.Errorf("Update should have returned err")
			}
		})
	}
}

func TestVault_Update_ClosedServer(t *testing.T) {
	// setup types
	sec := new(library.Secret)
	sec.SetOrg("foo")
	sec.SetRepo("bar")
	sec.SetName("baz")
	sec.SetValue("foob")
	sec.SetType("repo")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetEvents([]string{"foo", "bar"})

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
			s, err := New(fake.URL, "foo", tt.args.version, tt.args.prefix)
			if err != nil {
				t.Errorf("New returned err: %v", err)
			}

			err = s.Update("repo", "foo", "bar", sec)
			if err == nil {
				t.Errorf("Update should have returned err")
			}
		})
	}
}

func TestVault_Update_NoWrite(t *testing.T) {
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
	engine.PUT("/v1/secret/repo/foo/bar/baz", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	engine.GET("/v1/secret/data/repo/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/repo.json")
	})
	engine.PUT("/v1/secret/data/repo/foo/bar/baz", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	engine.GET("/v1/secret/data/prefix/repo/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/repo.json")
	})
	engine.PUT("/v1/secret/data/prefix/repo/foo/bar/baz", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
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
			s, err := New(fake.URL, "foo", tt.args.version, tt.args.prefix)
			if err != nil {
				t.Errorf("New returned err: %v", err)
			}

			err = s.Update("repo", "foo", "bar", sec)

			if resp.Code != http.StatusOK {
				t.Errorf("Update returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err == nil {
				t.Errorf("Update should have returned err")
			}
		})
	}
}
