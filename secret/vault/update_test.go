// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	api "github.com/go-vela/server/api/types"
)

func TestVault_Update_Org(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.PUT("/v1/secret/org/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/org.json")
	})
	engine.GET("/v1/secret/org/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/org.json")
	})

	engine.PUT("/v1/secret/data/org/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/org.json")
	})
	engine.GET("/v1/secret/data/org/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/org.json")
	})

	engine.PUT("/v1/secret/data/prefix/org/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/org.json")
	})
	engine.GET("/v1/secret/data/prefix/org/foo/bar", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/org.json")
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	// setup types
	sec := new(api.Secret)
	sec.SetOrg("foo")
	sec.SetRepo("*")
	sec.SetName("bar")
	sec.SetValue("baz")
	sec.SetType("org")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetAllowCommand(true)
	sec.SetAllowSubstitution(true)
	sec.SetAllowEvents(api.NewEventsFromMask(1))
	sec.SetCreatedAt(1563474077)
	sec.SetCreatedBy("octocat")
	sec.SetUpdatedAt(1563474079)
	sec.SetUpdatedBy("octocat2")

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

			got, err := s.Update(context.TODO(), "org", "foo", "*", sec)

			if resp.Code != http.StatusOK {
				t.Errorf("Update returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err != nil {
				t.Errorf("Update returned err: %v", err)
			}

			if !reflect.DeepEqual(got, sec) {
				t.Errorf("Update returned %s, want %s", got, sec)
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
	engine.PUT("/v1/secret/repo/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/repo.json")
	})
	engine.GET("/v1/secret/repo/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/repo.json")
	})

	engine.PUT("/v1/secret/data/repo/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/repo.json")
	})
	engine.GET("/v1/secret/data/repo/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/repo.json")
	})

	engine.PUT("/v1/secret/data/prefix/repo/foo/bar/baz", func(c *gin.Context) {
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
	sec := new(api.Secret)
	sec.SetOrg("foo")
	sec.SetRepo("bar")
	sec.SetName("baz")
	sec.SetValue("foob")
	sec.SetType("repo")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetAllowCommand(true)
	sec.SetAllowSubstitution(true)
	sec.SetAllowEvents(api.NewEventsFromMask(3))
	sec.SetCreatedAt(1563474077)
	sec.SetCreatedBy("octocat")
	sec.SetUpdatedAt(1563474079)
	sec.SetUpdatedBy("octocat2")

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

			got, err := s.Update(context.TODO(), "repo", "foo", "bar", sec)

			if resp.Code != http.StatusOK {
				t.Errorf("Update returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err != nil {
				t.Errorf("Update returned err: %v", err)
			}

			if !reflect.DeepEqual(got, sec) {
				t.Errorf("Update returned %s, want %s", got, sec)
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
	engine.PUT("/v1/secret/shared/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/shared.json")
	})
	engine.GET("/v1/secret/shared/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/shared.json")
	})

	engine.PUT("/v1/secret/data/shared/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/shared.json")
	})
	engine.GET("/v1/secret/data/shared/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/shared.json")
	})

	engine.PUT("/v1/secret/data/prefix/shared/foo/bar/baz", func(c *gin.Context) {
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
	sec := new(api.Secret)
	sec.SetOrg("foo")
	sec.SetTeam("bar")
	sec.SetName("baz")
	sec.SetValue("foob")
	sec.SetType("shared")
	sec.SetImages([]string{"foo", "bar"})
	sec.SetAllowCommand(false)
	sec.SetAllowSubstitution(false)
	sec.SetAllowEvents(api.NewEventsFromMask(1))
	sec.SetRepoAllowlist([]string{"github/octocat", "github/octokitty"})
	sec.SetCreatedAt(1563474077)
	sec.SetCreatedBy("octocat")
	sec.SetUpdatedAt(1563474079)
	sec.SetUpdatedBy("octocat2")

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

			got, err := s.Update(context.TODO(), "shared", "foo", "bar", sec)

			if resp.Code != http.StatusOK {
				t.Errorf("Update returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err != nil {
				t.Errorf("Update returned err: %v", err)
			}

			if !reflect.DeepEqual(got, sec) {
				t.Errorf("Update returned %s, want %s", got, sec)
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
	engine.PUT("/v1/secret/repo/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v1/invalid_repo.json")
	})

	engine.PUT("/v1/secret/data/repo/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/invalid_repo.json")
	})

	engine.PUT("/v1/secret/data/prefix/repo/foo/bar/baz", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/v2/invalid_repo.json")
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	// setup types
	sec := new(api.Secret)
	sec.SetOrg("foo")
	sec.SetRepo("bar")
	sec.SetName("baz")
	sec.SetValue("")
	sec.SetType("repo")
	sec.SetImages([]string{"foo", "bar"})
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

			_, err = s.Update(context.TODO(), "repo", "foo", "bar", sec)

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
	sec := new(api.Secret)
	sec.SetOrg("foo")
	sec.SetRepo("bar")
	sec.SetName("baz")
	sec.SetValue("foob")
	sec.SetType("invalid")
	sec.SetImages([]string{"foo", "bar"})

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

			_, err = s.Update(context.TODO(), "invalid", "foo", "bar", sec)
			if err == nil {
				t.Errorf("Update should have returned err")
			}
		})
	}
}

func TestVault_Update_ClosedServer(t *testing.T) {
	// setup types
	sec := new(api.Secret)
	sec.SetOrg("foo")
	sec.SetRepo("bar")
	sec.SetName("baz")
	sec.SetValue("foob")
	sec.SetType("repo")
	sec.SetImages([]string{"foo", "bar"})

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

			_, err = s.Update(context.TODO(), "repo", "foo", "bar", sec)
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
	engine.PUT("/v1/secret/repo/foo/bar/baz", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	engine.PUT("/v1/secret/data/repo/foo/bar/baz", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	engine.PUT("/v1/secret/data/prefix/repo/foo/bar/baz", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	fake := httptest.NewServer(engine)
	defer fake.Close()

	// setup types
	sec := new(api.Secret)
	sec.SetOrg("foo")
	sec.SetRepo("bar")
	sec.SetName("baz")
	sec.SetValue("foob")
	sec.SetType("repo")
	sec.SetImages([]string{"foo", "bar"})

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

			_, err = s.Update(context.TODO(), "repo", "foo", "bar", sec)

			if resp.Code != http.StatusOK {
				t.Errorf("Update returned %v, want %v", resp.Code, http.StatusOK)
			}

			if err == nil {
				t.Errorf("Update should have returned err")
			}
		})
	}
}
