// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package step

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-vela/server/database"

	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
)

func TestStep_Retrieve(t *testing.T) {
	// setup types
	want := new(library.Step)
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

func TestStep_Establish(t *testing.T) {
	// setup types
	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetHash("baz")
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")
	r.SetVisibility("public")

	b := new(library.Build)
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)

	want := new(library.Step)
	want.SetID(1)
	want.SetRepoID(1)
	want.SetBuildID(1)
	want.SetNumber(1)
	want.SetName("foo")
	want.SetImage("baz")
	want.SetStage("")
	want.SetStatus("")
	want.SetError("")
	want.SetExitCode(0)
	want.SetCreated(0)
	want.SetStarted(0)
	want.SetFinished(0)
	want.SetHost("")
	want.SetRuntime("")
	want.SetDistribution("")

	got := new(library.Step)

	// setup database
	db, _ := database.NewTest()

	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Exec("delete from builds;")
		db.Database.Exec("delete from steps;")
		db.Database.Close()
	}()

	_ = db.CreateRepo(r)
	_ = db.CreateBuild(b)
	_ = db.CreateStep(want)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/1/steps/1", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(repo.Establish())
	engine.Use(build.Establish())
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds/:build/steps/:step", func(c *gin.Context) {
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

func TestStep_Establish_NoRepo(t *testing.T) {
	// setup database
	db, _ := database.NewTest()
	defer db.Database.Close()

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/1/steps/1", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds/:build/steps/:step", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusNotFound)
	}
}

func TestStep_Establish_NoBuild(t *testing.T) {
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
	db, _ := database.NewTest()

	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Close()
	}()

	_ = db.CreateRepo(r)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/1/steps/1", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(repo.Establish())
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds/:build/steps/:step", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusNotFound)
	}
}

func TestStep_Establish_NoStepParameter(t *testing.T) {
	// setup types
	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetHash("baz")
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")
	r.SetVisibility("public")

	b := new(library.Build)
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)

	// setup database
	db, _ := database.NewTest()

	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Exec("delete from builds;")
		db.Database.Close()
	}()

	_ = db.CreateRepo(r)
	_ = db.CreateBuild(b)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/1/steps", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(repo.Establish())
	engine.Use(build.Establish())
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds/:build/steps", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusBadRequest)
	}
}

func TestStep_Establish_InvalidStepParameter(t *testing.T) {
	// setup types
	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetHash("baz")
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")
	r.SetVisibility("public")

	b := new(library.Build)
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)

	// setup database
	db, _ := database.NewTest()

	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Exec("delete from builds;")
		db.Database.Close()
	}()

	_ = db.CreateRepo(r)
	_ = db.CreateBuild(b)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/1/steps/foo", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(repo.Establish())
	engine.Use(build.Establish())
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds/:build/steps/:step", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusBadRequest)
	}
}

func TestStep_Establish_NoStep(t *testing.T) {
	// setup types
	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetHash("baz")
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")
	r.SetVisibility("public")

	b := new(library.Build)
	b.SetID(1)
	b.SetRepoID(1)
	b.SetNumber(1)

	// setup database
	db, _ := database.NewTest()

	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Exec("delete from builds;")
		db.Database.Close()
	}()

	_ = db.CreateRepo(r)
	_ = db.CreateBuild(b)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/1/steps/1", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(repo.Establish())
	engine.Use(build.Establish())
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds/:build/steps/:step", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusNotFound)
	}
}
