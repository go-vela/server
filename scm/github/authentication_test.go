// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	_context "context"

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
	got, err := client.Authenticate(_context.TODO(), context.Writer, context.Request, "bar")

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
	got, err := client.Authenticate(_context.TODO(), context.Writer, context.Request, "bar")

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
	got, err := client.Authenticate(_context.TODO(), context.Writer, context.Request, "bar")

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
	got, err := client.Authenticate(_context.TODO(), context.Writer, context.Request, "bar")

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
	got, err := client.Authenticate(_context.TODO(), context.Writer, context.Request, "bar")

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
	got, err := client.Authorize(_context.TODO(), "foobar")

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
	got, err := client.Authorize(_context.TODO(), "foobar")

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
	_, err := client.Login(_context.TODO(), context.Writer, context.Request)

	if resp.Code != http.StatusTemporaryRedirect {
		t.Errorf("Login returned %v, want %v", resp.Code, http.StatusTemporaryRedirect)
	}

	if err != nil {
		t.Errorf("Login returned err: %v", err)
	}
}

func TestGithub_AuthenticateToken(t *testing.T) {
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
	got, err := client.AuthenticateToken(_context.TODO(), context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("AuthenticateToken returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("AuthenticateToken returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("AuthenticateToken is %v, want %v", got, want)
	}
}

func TestGithub_AuthenticateToken_Invalid(t *testing.T) {
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
	got, err := client.AuthenticateToken(_context.TODO(), context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("AuthenticateToken returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Errorf("AuthenticateToken did not return err")
	}

	if got != nil {
		t.Errorf("AuthenticateToken is %v, want nil", got)
	}
}

func TestGithub_AuthenticateToken_Vela_OAuth(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodPost, "/authenticate/token", nil)
	context.Request.Header.Set("Token", "vela")

	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/user.json")
	})

	engine.POST("/api/v3/applications/foo/token", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	client, _ := NewTest(s.URL)

	// run test
	_, err := client.AuthenticateToken(_context.TODO(), context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("AuthenticateToken returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err == nil {
		t.Error("AuthenticateToken should have returned err")
	}
}

func TestGithub_ValidateOAuthToken_Valid(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/validate-oauth", nil)

	token := "foobar"
	want := true
	scmResponseCode := http.StatusOK

	engine.POST("/api/v3/applications/foo/token", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(scmResponseCode)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.ValidateOAuthToken(_context.TODO(), token)

	if got != want {
		t.Errorf("ValidateOAuthToken returned %v, want %v", got, want)
	}

	if err != nil {
		t.Errorf("ValidateOAuthToken returned err: %v", err)
	}
}

func TestGithub_ValidateOAuthToken_Invalid(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/validate-oauth", nil)

	token := "foobar"
	want := false
	// 404 from the mocked github server indicates an invalid oauth token
	scmResponseCode := http.StatusNotFound

	engine.POST("/api/v3/applications/foo/token", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(scmResponseCode)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.ValidateOAuthToken(_context.TODO(), token)

	if got != want {
		t.Errorf("ValidateOAuthToken returned %v, want %v", got, want)
	}

	if err != nil {
		t.Errorf("ValidateOAuthToken returned err: %v", err)
	}
}

func TestGithub_ValidateOAuthToken_Error(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/validate-oauth", nil)

	token := "foobar"
	want := false
	scmResponseCode := http.StatusInternalServerError

	engine.POST("/api/v3/applications/foo/token", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(scmResponseCode)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	client, _ := NewTest(s.URL)

	// run test
	got, err := client.ValidateOAuthToken(_context.TODO(), token)

	if got != want {
		t.Errorf("ValidateOAuthToken returned %v, want %v", got, want)
	}

	if err == nil {
		t.Errorf("ValidateOAuthToken did not return err")
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
	_, err := client.Login(_context.TODO(), context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Login returned %v, want %v", resp.Code, http.StatusOK)
	}

	if err != nil {
		t.Errorf("Login returned err: %v", err)
	}
}
