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

	want := new(api.Secret)
	want.SetOrg("foo")
	want.SetRepo("*")
	want.SetName("bar")
	want.SetValue("baz")
	want.SetType("org")
	want.SetImages([]string{"foo", "bar"})
	want.SetAllowCommand(true)
	want.SetAllowSubstitution(true)
	want.SetAllowEvents(api.NewEventsFromMask(1))
	want.SetCreatedAt(1563474077)
	want.SetCreatedBy("octocat")
	want.SetUpdatedAt(1563474079)
	want.SetUpdatedBy("octocat2")

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
			got, err := s.Get(context.TODO(), "org", "foo", "bar", "baz")

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
	want := new(api.Secret)
	want.SetOrg("foo")
	want.SetRepo("bar")
	want.SetName("baz")
	want.SetValue("foob")
	want.SetType("repo")
	want.SetImages([]string{"foo", "bar"})
	want.SetAllowCommand(true)
	want.SetAllowSubstitution(true)
	want.SetAllowEvents(api.NewEventsFromMask(3))
	want.SetCreatedAt(1563474077)
	want.SetCreatedBy("octocat")
	want.SetUpdatedAt(1563474079)
	want.SetUpdatedBy("octocat2")

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
			got, err := s.Get(context.TODO(), "repo", "foo", "bar", "baz")

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
	want := new(api.Secret)
	want.SetOrg("foo")
	want.SetTeam("bar")
	want.SetName("baz")
	want.SetValue("foob")
	want.SetType("shared")
	want.SetImages([]string{"foo", "bar"})
	want.SetAllowCommand(false)
	want.SetAllowSubstitution(false)
	want.SetRepoAllowlist([]string{"github/octocat", "github/octokitty"})
	want.SetAllowEvents(api.NewEventsFromMask(1))
	want.SetCreatedAt(1563474077)
	want.SetCreatedBy("octocat")
	want.SetUpdatedAt(1563474079)
	want.SetUpdatedBy("octocat2")

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
			got, err := s.Get(context.TODO(), "shared", "foo", "bar", "baz")

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
			got, err := s.Get(context.TODO(), "invalid", "foo", "bar", "foob")
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
			got, err := s.Get(context.TODO(), "repo", "foo", "bar", "foob")
			if err == nil {
				t.Errorf("Get should have returned err")
			}

			if got != nil {
				t.Errorf("Get is %v, want nil", got)
			}
		})
	}
}
