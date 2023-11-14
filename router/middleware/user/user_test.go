// SPDX-License-Identifier: Apache-2.0

package user

import (
	_context "context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/scm/github"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/golang-jwt/jwt/v5"
)

func TestUser_Retrieve(t *testing.T) {
	// setup types
	want := new(library.User)
	want.SetID(1)

	// setup context
	gin.SetMode(gin.TestMode)

	context, _ := gin.CreateTestContext(nil)
	ToContext(context, want)

	// run test
	got := Retrieve(context)

	if got != want {
		t.Errorf("Retrieve is %v, want %v", got, want)
	}
}

func TestUser_Establish(t *testing.T) {
	// setup types
	secret := "superSecret"

	tm := &token.Manager{
		PrivateKey:               "123abc",
		SignMethod:               jwt.SigningMethodHS256,
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}

	want := new(library.User)
	want.SetID(1)
	want.SetName("foo")
	want.SetRefreshToken("fresh")
	want.SetToken("bar")
	want.SetHash("baz")
	want.SetActive(false)
	want.SetAdmin(false)
	want.SetFavorites([]string{})

	got := new(library.User)

	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/users/foo", nil)

	mto := &token.MintTokenOpts{
		User:          want,
		TokenDuration: tm.UserAccessTokenDuration,
		TokenType:     constants.UserAccessTokenType,
	}

	at, _ := tm.MintToken(mto)

	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", at))
	context.Request.AddCookie(&http.Cookie{
		Name:  constants.RefreshTokenName,
		Value: "fresh",
	})

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteUser(_context.TODO(), want)
		db.Close()
	}()

	_, _ = db.CreateUser(_context.TODO(), want)

	// setup context
	gin.SetMode(gin.TestMode)

	// setup github mock server
	engine.GET("/api/v3/user", func(c *gin.Context) {
		c.String(http.StatusOK, userPayload)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := github.NewTest(s.URL)

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { scm.ToContext(c, client) })
	engine.Use(claims.Establish())
	engine.Use(Establish())
	engine.GET("/users/:user", func(c *gin.Context) {
		got = Retrieve(c)

		c.Status(http.StatusOK)
	})

	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Establish is %v, want %v", got, want)
	}
}

func TestUser_Establish_NoToken(t *testing.T) {
	// setup types
	secret := "superSecret"
	tm := &token.Manager{
		PrivateKey:               "123abc",
		SignMethod:               jwt.SigningMethodHS256,
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}
	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}
	defer db.Close()

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/users/foo", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(claims.Establish())
	engine.Use(Establish())

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusUnauthorized)
	}
}

func TestUser_Establish_DiffTokenType(t *testing.T) {
	// setup types
	secret := "superSecret"

	tm := &token.Manager{
		PrivateKey:               "123abc",
		SignMethod:               jwt.SigningMethodHS256,
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}

	want := new(library.User)

	got := new(library.User)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/users/foo", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", secret))

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(claims.Establish())
	engine.Use(Establish())
	engine.GET("/users/:user", func(c *gin.Context) {
		got = Retrieve(c)

		c.Status(http.StatusOK)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Establish is %v, want %v", got, want)
	}
}

func TestUser_Establish_NoAuthorizeUser(t *testing.T) {
	// setup database
	secret := "superSecret"

	tm := &token.Manager{
		PrivateKey:               "123abc",
		SignMethod:               jwt.SigningMethodHS256,
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}
	defer db.Close()

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/users/foo?access_token=bar", nil)

	// setup client
	client, _ := github.NewTest("")

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { scm.ToContext(c, client) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(claims.Establish())
	engine.Use(Establish())

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusUnauthorized)
	}
}

func TestUser_Establish_NoUser(t *testing.T) {
	// setup types
	tm := &token.Manager{
		PrivateKey:               "123abc",
		SignMethod:               jwt.SigningMethodHS256,
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}

	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")

	// setup database
	secret := "superSecret"

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}
	defer db.Close()

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/users/foo?access_token=bar", nil)

	mto := &token.MintTokenOpts{
		User:          u,
		TokenDuration: tm.UserAccessTokenDuration,
		TokenType:     constants.UserAccessTokenType,
	}

	at, _ := tm.MintToken(mto)

	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", at))
	context.Request.AddCookie(&http.Cookie{
		Name:  constants.RefreshTokenName,
		Value: "fresh",
	})

	// setup client
	client, _ := github.NewTest("")

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { scm.ToContext(c, client) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(claims.Establish())
	engine.Use(Establish())

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusUnauthorized)
	}
}

const userPayload = `
{
  "login": "foo",
  "id": 1,
  "node_id": "MDQ6VXNlcjE=",
  "avatar_url": "https://github.com/images/error/octocat_happy.gif",
  "gravatar_id": "",
  "url": "https://api.github.com/users/foo",
  "html_url": "https://github.com/octocat",
  "followers_url": "https://api.github.com/users/foo/followers",
  "following_url": "https://api.github.com/users/foo/following{/other_user}",
  "gists_url": "https://api.github.com/users/foo/gists{/gist_id}",
  "starred_url": "https://api.github.com/users/foo/starred{/org}{/repo}",
  "subscriptions_url": "https://api.github.com/users/foo/subscriptions",
  "orgs_url": "https://api.github.com/users/foo/orgs",
  "repos_url": "https://api.github.com/users/foo/repos",
  "events_url": "https://api.github.com/users/foo/events{/privacy}",
  "received_events_url": "https://api.github.com/users/foo/received_events",
  "type": "User",
  "site_admin": false,
  "name": "monalisa foo",
  "company": "GitHub",
  "blog": "https://github.com/blog",
  "location": "San Francisco",
  "email": "foo@github.com",
  "hireable": false,
  "bio": "There once was...",
  "public_repos": 2,
  "public_gists": 1,
  "followers": 20,
  "following": 0,
  "created_at": "2008-01-14T04:33:35Z",
  "updated_at": "2008-01-14T04:33:35Z",
  "private_gists": 81,
  "total_private_repos": 100,
  "owned_private_repos": 100,
  "disk_usage": 10000,
  "collaborators": 8,
  "two_factor_authentication": true,
  "plan": {
    "name": "Medium",
    "space": 400,
    "private_repos": 20,
    "collaborators": 0
  }
}
`
