// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package perm

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/token"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/source"
	"github.com/go-vela/server/source/github"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
)

func TestPerm_MustPlatformAdmin(t *testing.T) {
	// setup types
	secret := "superSecret"

	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")
	u.SetAdmin(true)

	tkn, err := token.Compose(u)
	if err != nil {
		t.Errorf("Unable to Compose token: %v", err)
	}

	// setup database
	db, _ := database.NewTest()
	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()
	_ = db.CreateUser(u)

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/admin/users", nil)
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
	engine.Use(user.Establish())
	engine.Use(MustPlatformAdmin())
	engine.GET("/admin/users", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustPlatAdmin returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustPlatformAdmin_NotAdmin(t *testing.T) {
	// setup types
	secret := "superSecret"

	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")
	u.SetAdmin(false)

	tkn, err := token.Compose(u)
	if err != nil {
		t.Errorf("Unable to Compose token: %v", err)
	}

	// setup database
	db, _ := database.NewTest()
	defer func() {
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()
	_ = db.CreateUser(u)

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/admin/users", nil)
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
	engine.Use(user.Establish())
	engine.Use(MustPlatformAdmin())
	engine.GET("/admin/users", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("MustPlatAdmin returned %v, want %v", resp.Code, http.StatusUnauthorized)
	}
}

func TestPerm_MustAdmin(t *testing.T) {
	// setup types
	secret := "superSecret"

	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")
	u.SetAdmin(false)

	tkn, err := token.Compose(u)
	if err != nil {
		t.Errorf("Unable to Compose token: %v", err)
	}

	// setup database
	db, _ := database.NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(r)
	_ = db.CreateUser(u)

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permAdminPayload)
	})
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
	engine.Use(user.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustAdmin())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustAdmin returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustAdmin_PlatAdmin(t *testing.T) {
	// setup types
	secret := "superSecret"

	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")
	u.SetAdmin(true)

	tkn, err := token.Compose(u)
	if err != nil {
		t.Errorf("Unable to Compose token: %v", err)
	}

	// setup database
	db, _ := database.NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(r)
	_ = db.CreateUser(u)

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permWritePayload)
	})
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
	engine.Use(user.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustAdmin())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustAdmin returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustAdmin_NotAdmin(t *testing.T) {
	// setup types
	secret := "superSecret"

	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")
	u.SetAdmin(false)

	tkn, err := token.Compose(u)
	if err != nil {
		t.Errorf("Unable to Compose token: %v", err)
	}

	// setup database
	db, _ := database.NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(r)
	_ = db.CreateUser(u)

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permWritePayload)
	})
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
	engine.Use(user.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustAdmin())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("MustAdmin returned %v, want %v", resp.Code, http.StatusUnauthorized)
	}
}

func TestPerm_MustWrite(t *testing.T) {
	// setup types
	secret := "superSecret"

	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")
	u.SetAdmin(false)

	tkn, err := token.Compose(u)
	if err != nil {
		t.Errorf("Unable to Compose token: %v", err)
	}

	// setup database
	db, _ := database.NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(r)
	_ = db.CreateUser(u)

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permWritePayload)
	})
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
	engine.Use(user.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustWrite())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustWrite returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustWrite_PlatAdmin(t *testing.T) {
	// setup types
	secret := "superSecret"

	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")
	u.SetAdmin(true)

	tkn, err := token.Compose(u)
	if err != nil {
		t.Errorf("Unable to Compose token: %v", err)
	}

	// setup database
	db, _ := database.NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(r)
	_ = db.CreateUser(u)

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permWritePayload)
	})
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
	engine.Use(user.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustWrite())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustWrite returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustWrite_RepoAdmin(t *testing.T) {
	// setup types
	secret := "superSecret"

	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")
	u.SetAdmin(false)

	tkn, err := token.Compose(u)
	if err != nil {
		t.Errorf("Unable to Compose token: %v", err)
	}

	// setup database
	db, _ := database.NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(r)
	_ = db.CreateUser(u)

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permAdminPayload)
	})
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
	engine.Use(user.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustWrite())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustWrite returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustWrite_NotWrite(t *testing.T) {
	// setup types
	secret := "superSecret"

	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")
	u.SetAdmin(false)

	tkn, err := token.Compose(u)
	if err != nil {
		t.Errorf("Unable to Compose token: %v", err)
	}

	// setup database
	db, _ := database.NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(r)
	_ = db.CreateUser(u)

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permReadPayload)
	})
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
	engine.Use(user.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustWrite())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("MustWrite returned %v, want %v", resp.Code, http.StatusUnauthorized)
	}
}

func TestPerm_MustRead(t *testing.T) {
	// setup types
	secret := "superSecret"

	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")
	u.SetAdmin(false)

	tkn, err := token.Compose(u)
	if err != nil {
		t.Errorf("Unable to Compose token: %v", err)
	}

	// setup database
	db, _ := database.NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(r)
	_ = db.CreateUser(u)

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permReadPayload)
	})
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
	engine.Use(user.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustRead())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustRead returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustRead_PlatAdmin(t *testing.T) {
	// setup types
	secret := "superSecret"

	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")
	u.SetAdmin(true)

	tkn, err := token.Compose(u)
	if err != nil {
		t.Errorf("Unable to Compose token: %v", err)
	}

	// setup database
	db, _ := database.NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(r)
	_ = db.CreateUser(u)

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permReadPayload)
	})
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
	engine.Use(user.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustRead())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustRead returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustRead_RepoAdmin(t *testing.T) {
	// setup types
	secret := "superSecret"

	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")
	u.SetAdmin(false)

	tkn, err := token.Compose(u)
	if err != nil {
		t.Errorf("Unable to Compose token: %v", err)
	}

	// setup database
	db, _ := database.NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(r)
	_ = db.CreateUser(u)

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permAdminPayload)
	})
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
	engine.Use(user.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustRead())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustRead returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustRead_RepoWrite(t *testing.T) {
	// setup types
	secret := "superSecret"

	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")
	u.SetAdmin(false)

	tkn, err := token.Compose(u)
	if err != nil {
		t.Errorf("Unable to Compose token: %v", err)
	}

	// setup database
	db, _ := database.NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(r)
	_ = db.CreateUser(u)

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permWritePayload)
	})
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
	engine.Use(user.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustRead())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("MustRead returned %v, want %v", resp.Code, http.StatusOK)
	}
}

func TestPerm_MustRead_NotRead(t *testing.T) {
	// setup types
	secret := "superSecret"

	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")

	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")
	u.SetAdmin(false)

	tkn, err := token.Compose(u)
	if err != nil {
		t.Errorf("Unable to Compose token: %v", err)
	}

	// setup database
	db, _ := database.NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Exec("delete from users;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(r)
	_ = db.CreateUser(u)

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/test/foo/bar", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tkn))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/collaborators/:username/permission", func(c *gin.Context) {
		c.String(http.StatusOK, permNonePayload)
	})
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
	engine.Use(user.Establish())
	engine.Use(repo.Establish())
	engine.Use(MustRead())
	engine.GET("/test/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	s1 := httptest.NewServer(engine)
	defer s1.Close()

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("MustRead returned %v, want %v", resp.Code, http.StatusUnauthorized)
	}
}

func TestPerm_globalPerms(t *testing.T) {
	// setup types
	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")
	u.SetAdmin(false)

	// run test
	got := globalPerms(u)

	if got {
		t.Errorf("globalPerms returned %v, want false", got)
	}
}

func TestPerm_globalPerms_Agent(t *testing.T) {
	// setup types
	u := new(library.User)
	u.SetID(1)
	u.SetName("vela-worker")
	u.SetToken("bar")
	u.SetHash("baz")
	u.SetAdmin(false)

	// run test
	got := globalPerms(u)

	if !got {
		t.Errorf("globalPerms returned %v, want true", got)
	}
}

func TestPerm_globalPerms_Admin(t *testing.T) {
	// setup types
	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")
	u.SetAdmin(true)

	// run test
	got := globalPerms(u)

	if !got {
		t.Errorf("globalPerms returned %v, want true", got)
	}
}

const permAdminPayload = `
{
  "permission": "admin",
  "user": {
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
    "site_admin": false
  }
}
`

const permWritePayload = `
{
  "permission": "write",
  "user": {
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
    "site_admin": false
  }
}
`

const permReadPayload = `
{
  "permission": "read",
  "user": {
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
    "site_admin": false
  }
}
`

const permNonePayload = `
{
  "permission": "none",
  "user": {
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
    "site_admin": false
  }
}
`

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
