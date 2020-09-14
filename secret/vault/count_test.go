// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestVault_Count_Org(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/v1/secret/org/foo/foo", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/org.json")
	})
	engine.GET("/v1/secret/org/foo", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/list.json")
	})

	engine.GET("/v1/secret/data/org/foo/foo", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/org.json")
	})
	engine.GET("/v1/secret/metadata/org/foo", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/list.json")
	})

	engine.GET("/v1/secret/data/prefix/org/foo/foo", func(c *gin.Context) {
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
	want := 1

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

			got, err := s.Count("org", "foo", "*")

			if resp.Code != http.StatusOK {
				t.Errorf("Count returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err != nil {
				t.Errorf("Count returned err: %v", err)
			}

			if got != int64(want) {
				t.Errorf("Count is %v, want %v", got, want)
			}
		})
	}
}

func TestVault_Count_Repo(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/v1/secret/repo/foo/bar/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/repo.json")
	})
	engine.GET("/v1/secret/repo/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/list.json")
	})

	engine.GET("/v1/secret/data/repo/foo/bar/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/repo.json")
	})
	engine.GET("/v1/secret/metadata/repo/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/list.json")
	})

	engine.GET("/v1/secret/data/prefix/repo/foo/bar/bar", func(c *gin.Context) {
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
	want := 1

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

			got, err := s.Count("repo", "foo", "bar")

			if resp.Code != http.StatusOK {
				t.Errorf("Count returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err != nil {
				t.Errorf("Count returned err: %v", err)
			}

			if got != int64(want) {
				t.Errorf("Count is %v, want %v", got, want)
			}
		})
	}
}

func TestVault_Count_Shared(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/v1/secret/shared/foo/bar/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/shared.json")
	})
	engine.GET("/v1/secret/shared/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/list.json")
	})

	engine.GET("/v1/secret/data/shared/foo/bar/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/shared.json")
	})
	engine.GET("/v1/secret/metadata/shared/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/list.json")
	})

	engine.GET("/v1/secret/data/prefix/shared/foo/bar/bar", func(c *gin.Context) {
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
	want := 1

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

			got, err := s.Count("shared", "foo", "bar")

			if resp.Code != http.StatusOK {
				t.Errorf("List returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err != nil {
				t.Errorf("List returned err: %v", err)
			}

			if got != int64(want) {
				t.Errorf("Count is %v, want %v", got, want)
			}
		})
	}
}

func TestVault_Count_InvalidType(t *testing.T) {
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

			got, err := s.Count("invalid", "foo", "bar")
			if err == nil {
				t.Errorf("Count should have returned err")
			}

			if got != 0 {
				t.Errorf("Count is %v, want 0", got)
			}
		})
	}
}

func TestVault_Count_ClosedServer(t *testing.T) {
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

			got, err := s.Count("repo", "foo", "bar")
			if err == nil {
				t.Errorf("Count should have returned err")
			}

			if got != 0 {
				t.Errorf("Count is %v, want 0", got)
			}
		})
	}
}

func TestVault_Count_EmptyList(t *testing.T) {
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
			s, err := New(fake.URL, "foo", tt.args.version, tt.args.prefix)
			if err != nil {
				t.Errorf("New returned err: %v", err)
			}

			got, err := s.Count("repo", "foo", "bar")

			if resp.Code != http.StatusOK {
				t.Errorf("Count returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err == nil {
				t.Errorf("Count should have returned err")
			}

			if got != 0 {
				t.Errorf("Count is %v, want 0", got)
			}
		})
	}
}

func TestVault_Count_InvalidList(t *testing.T) {
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
			s, err := New(fake.URL, "foo", tt.args.version, tt.args.prefix)
			if err != nil {
				t.Errorf("New returned err: %v", err)
			}

			got, err := s.Count("repo", "foo", "bar")

			if resp.Code != http.StatusOK {
				t.Errorf("Count returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err == nil {
				t.Errorf("Count should have returned err")
			}

			if got != 0 {
				t.Errorf("Count is %v, want 0", got)
			}
		})
	}
}
