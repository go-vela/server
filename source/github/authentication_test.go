// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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

func TestGithub_Authenticate(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/login/oauth/authorize?code=foo&state=bar", nil)

	// setup mock server
	engine.POST("/login/oauth/access_token", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/token.json")
	})
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/user.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	want := new(library.User)
	want.SetName("octocat")
	want.SetToken("foo")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Authenticate(context.Writer, context.Request, "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Authenticate returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Authenticate returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Authenticate is %v, want %v", got, want)
	}
}

func TestGithub_Authenticate_NoCode(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/login", nil)

	// setup mock server
	engine.Any("/", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Authenticate(context.Writer, context.Request, "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Authenticate returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Authenticate returned err: %v", err)
	}

	if got != nil {
		t.Errorf("Authenticate is %v, want nil", got)
	}
}

func TestGithub_Authenticate_NoState(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/login?code=foo", nil)

	// setup mock server
	engine.Any("/", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Authenticate(context.Writer, context.Request, "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Authenticate returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Errorf("Authenticate should have returned err")
	}

	if got != nil {
		t.Errorf("Authenticate is %v, want nil", got)
	}
}

func TestGithub_Authenticate_BadToken(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/login?code=foo&state=bar", nil)

	// setup mock server
	engine.Any("/", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Authenticate(context.Writer, context.Request, "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Authenticate returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Errorf("Authenticate should have returned err")
	}

	if got != nil {
		t.Errorf("Authenticate is %v, want nil", got)
	}
}

func TestGithub_Authenticate_NotFound(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/login/oauth/authorize?code=foo&state=bar", nil)

	// setup mock server
	engine.POST("/login/oauth/access_token", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/token.json")
	})
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Authenticate(context.Writer, context.Request, "bar")

	if resp.Code != http.StatusOK {
		t.Errorf("Authenticate returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Errorf("Authenticate should have returned err")
	}

	if got != nil {
		t.Errorf("Authenticate is %v, want nil", got)
	}
}

func TestGithub_Authorize(t *testing.T) {
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

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	want := "octocat"
	got, err := client.Authorize("foobar")

	if resp.Code != http.StatusOK {
		t.Errorf("Authorize returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Authorize returned err: %v", err)
	}

	if got != want {
		t.Errorf("Authorize is %v, want %v", got, want)
	}
}

func TestGithub_Authorize_NotFound(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	got, err := client.Authorize("foobar")

	if resp.Code != http.StatusOK {
		t.Errorf("Authorize returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Errorf("Authorize should have returned err")
	}

	if len(got) > 0 {
		t.Errorf("Authorize is %v, want \"\"", got)
	}
}

func TestGithub_Login(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/login", nil)

	// setup mock server
	engine.Any("/", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	client, _ := NewTest(s.URL)

	// run test
	_, err := client.Login(context.Writer, context.Request)

	if resp.Code != http.StatusTemporaryRedirect {
		t.Errorf("Login returned %v, want %v", resp.Code, http.StatusTemporaryRedirect)
	}

	if err != nil {
		t.Errorf("Login returned err: %v", err)
	}
}

func TestGithub_Authenticate_Token(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodPost, "/authenticate/token", nil)
	context.Request.Header.Set("Token", "foo")

	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/user.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	want := new(library.User)
	want.SetName("octocat")
	want.SetToken("foo")

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.AuthenticateToken(context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Authenticate returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Authenticate returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Authenticate is %v, want %v", got, want)
	}
}

func TestGithub_Authenticate_Invalid_Token(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodPost, "/authenticate/token", nil)
	context.Request.Header.Set("Token", "foo")

	// setup mock server
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := NewTest(s.URL)

	// run test
	got, err := client.AuthenticateToken(context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Authenticate returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Errorf("Authenticate did not return err")
	}

	if got != nil {
		t.Errorf("Authenticate is %v, want nil", got)
	}
}

func TestGithub_LoginWCreds(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodPost, "/login", nil)

	// setup mock server
	engine.Any("/", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	client, _ := NewTest(s.URL)

	// run test
	_, err := client.Login(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Enable returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Login returned err: %v", err)
	}
}
