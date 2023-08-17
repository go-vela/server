// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/compiler/native"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal/token"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/scm/github"
	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/golang-jwt/jwt/v5"
	"github.com/urfave/cli/v2"
)

func TestPipeline_Retrieve(t *testing.T) {
	// setup types
	_pipeline := new(library.Pipeline)

	gin.SetMode(gin.TestMode)
	_context, _ := gin.CreateTestContext(nil)

	// setup tests
	tests := []struct {
		name    string
		context *gin.Context
		want    *library.Pipeline
	}{
		{
			name:    "context",
			context: _context,
			want:    _pipeline,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ToContext(test.context, test.want)

			got := Retrieve(test.context)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Retrieve for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}

func TestPipeline_Establish(t *testing.T) {
	// setup types
	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetHash("baz")
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")
	r.SetVisibility("public")

	want := new(library.Pipeline)
	want.SetID(1)
	want.SetRepoID(1)
	want.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	want.SetFlavor("")
	want.SetPlatform("")
	want.SetRef("refs/heads/master")
	want.SetType("yaml")
	want.SetVersion("1")
	want.SetExternalSecrets(false)
	want.SetInternalSecrets(false)
	want.SetServices(false)
	want.SetStages(false)
	want.SetSteps(false)
	want.SetTemplates(false)
	want.SetData([]byte{})

	got := new(library.Pipeline)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		db.DeletePipeline(want)
		db.DeleteRepo(context.TODO(), r)
		db.Close()
	}()

	_, _ = db.CreateRepo(context.TODO(), r)
	_, _ = db.CreatePipeline(want)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/pipelines/foo/bar/48afb5bdc41ad69bf22588491333f7cf71135163", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(Establish())
	engine.GET("/pipelines/:org/:repo/:pipeline", func(c *gin.Context) {
		got = Retrieve(c)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Establish is %v, want %v", got, want)
	}
}

func TestPipeline_Establish_NoRepo(t *testing.T) {
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
	context.Request, _ = http.NewRequest(http.MethodGet, "/pipelines/foo/bar/48afb5bdc41ad69bf22588491333f7cf71135163", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(Establish())

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusNotFound)
	}
}

func TestPipeline_Establish_NoPipelineParameter(t *testing.T) {
	// setup types
	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetHash("baz")
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")
	r.SetVisibility("public")

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		db.DeleteRepo(context.TODO(), r)
		db.Close()
	}()

	_, _ = db.CreateRepo(context.TODO(), r)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/pipelines/foo/bar", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(Establish())
	engine.GET("/pipelines/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusBadRequest)
	}
}

func TestPipeline_Establish_NoPipeline(t *testing.T) {
	// setup types
	secret := "superSecret"

	tm := &token.Manager{
		PrivateKey:               "123abc",
		SignMethod:               jwt.SigningMethodHS256,
		UserAccessTokenDuration:  time.Minute * 5,
		UserRefreshTokenDuration: time.Minute * 30,
	}

	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetHash("baz")
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")
	r.SetVisibility("public")

	u := new(library.User)
	u.SetID(1)
	u.SetName("foo")
	u.SetToken("bar")
	u.SetHash("baz")
	u.SetAdmin(true)

	m := &types.Metadata{
		Database: &types.Database{
			Driver: "foo",
			Host:   "foo",
		},
		Queue: &types.Queue{
			Channel: "foo",
			Driver:  "foo",
			Host:    "foo",
		},
		Source: &types.Source{
			Driver: "foo",
			Host:   "foo",
		},
		Vela: &types.Vela{
			Address:    "foo",
			WebAddress: "foo",
		},
	}

	mto := &token.MintTokenOpts{
		User:          u,
		TokenDuration: tm.UserAccessTokenDuration,
		TokenType:     constants.UserAccessTokenType,
	}

	at, err := tm.MintToken(mto)
	if err != nil {
		t.Errorf("unable to mint user access token: %v", err)
	}

	set := flag.NewFlagSet("test", 0)
	set.String("clone-image", "target/vela-git:latest", "doc")

	comp, err := native.New(cli.NewContext(nil, set, nil))
	if err != nil {
		t.Errorf("unable to create compiler: %v", err)
	}

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		db.DeleteRepo(context.TODO(), r)
		db.DeleteUser(u)
		db.Close()
	}()

	_, _ = db.CreateRepo(context.TODO(), r)
	_ = db.CreateUser(u)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/pipelines/foo/bar/148afb5bdc41ad69bf22588491333f7cf71135163", nil)
	context.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", at))

	// setup github mock server
	engine.GET("/api/v3/repos/:org/:repo/contents/:path", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/yml.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup client
	client, _ := github.NewTest(s.URL)

	// setup vela mock server
	engine.Use(func(c *gin.Context) { c.Set("metadata", m) })
	engine.Use(func(c *gin.Context) { c.Set("token-manager", tm) })
	engine.Use(func(c *gin.Context) { c.Set("secret", secret) })
	engine.Use(func(c *gin.Context) { compiler.WithGinContext(c, comp) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(func(c *gin.Context) { scm.ToContext(c, client) })
	engine.Use(claims.Establish())
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(user.Establish())
	engine.Use(Establish())
	engine.GET("/pipelines/:org/:repo/:pipeline", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusOK)
	}
}
