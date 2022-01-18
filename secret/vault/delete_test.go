// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
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
	engine.DELETE("/v1/secret/org/foo/foob", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	engine.DELETE("/v1/secret/data/org/foo/foob", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	engine.DELETE("/v1/secret/data/prefix/org/foo/foob", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
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

			err = s.Delete("org", "foo", "bar", "foob")

			if resp.Code != http.StatusOK {
				t.Errorf("Delete returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err != nil {
				t.Errorf("Delete returned err: %v", err)
			}
		})
	}
}

func TestVault_Delete_Repo(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.DELETE("/v1/secret/repo/foo/bar/foob", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	engine.DELETE("/v1/secret/data/repo/foo/bar/foob", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	engine.DELETE("/v1/secret/data/prefix/repo/foo/bar/foob", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
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

			err = s.Delete("repo", "foo", "bar", "foob")

			if resp.Code != http.StatusOK {
				t.Errorf("Delete returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err != nil {
				t.Errorf("Delete returned err: %v", err)
			}
		})
	}
}

func TestVault_Delete_Shared(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.DELETE("/v1/secret/shared/foo/bar/foob", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	engine.DELETE("/v1/secret/data/shared/foo/bar/foob", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
	})

	engine.DELETE("/v1/secret/data/prefix/shared/foo/bar/foob", func(c *gin.Context) {
		c.String(http.StatusNoContent, "")
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

			err = s.Delete("shared", "foo", "bar", "foob")

			if resp.Code != http.StatusOK {
				t.Errorf("Delete returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err != nil {
				t.Errorf("Delete returned err: %v", err)
			}
		})
	}
}

func TestVault_Delete_InvalidType(t *testing.T) {
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

			err = s.Delete("invalid", "foo", "bar", "foob")
			if err == nil {
				t.Errorf("Delete should have returned err")
			}
		})
	}
}

func TestVault_Delete_ClosedServer(t *testing.T) {
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

			err = s.Delete("repo", "foo", "bar", "foob")
			if err == nil {
				t.Errorf("Delete should have returned err")
			}
		})
	}
}
