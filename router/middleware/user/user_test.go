// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package user

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/token"
	"github.com/go-vela/server/source"
	"github.com/go-vela/server/source/github"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
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

	want := new(library.User)
	want.SetID(1)
	want.SetName("foo")
	want.SetToken("bar")
	want.SetHash("baz")
	want.SetActive(false)
	want.SetAdmin(false)

	got := new(library.User)

	tkn, err := token.Compose(want)
	if err != nil {
		t.Errorf("Unable to Compose token: %v", err)
	}

	// setup database
	db, _ := database.NewTest()
	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()
	_ = db.CreateUser(want)

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/users/foo", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

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
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { source.ToContext(c, client) })
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
	// setup database
	db, _ := database.NewTest()
	defer db.Database.Close()

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/users/foo", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(Establish())

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusUnauthorized)
	}
}

func TestUser_Establish_SecretValid(t *testing.T) {
	// setup types
	secret := "superSecret"

	want := new(library.User)
	want.SetName("vela-worker")
	want.SetActive(true)
	want.SetAdmin(true)

	got := new(library.User)

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/users/foo", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", secret))

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
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
	db, _ := database.NewTest()
	defer db.Database.Close()

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/users/foo?access_token=bar", nil)

	// setup client
	client, _ := github.NewTest("")

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { source.ToContext(c, client) })
	engine.Use(Establish())

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusUnauthorized)
	}
}

func TestUser_Establish_NoUser(t *testing.T) {
	// setup types
	secret := "superSecret"
	got := new(library.User)

	// setup database
	db, _ := database.NewTest()
	defer db.Database.Close()

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/users/foo?access_token=bar", nil)

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
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { source.ToContext(c, client) })
	engine.Use(Establish())
	engine.GET("/users/:user", func(c *gin.Context) {
		got = Retrieve(c)

		c.Status(http.StatusOK)
	})
	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusUnauthorized)
	}

	if got.GetID() != 0 {
		t.Errorf("Establish is %v, want 0", got)
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
