// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package initstep

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-vela/server/router/middleware/org"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/database/sqlite"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/types/library"
)

func TestInitStep_Retrieve(t *testing.T) {
	// setup types
	want := new(library.InitStep)
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

func TestInitStep_Establish(t *testing.T) {
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

	want := new(library.InitStep)
	want.SetID(1)
	want.SetRepoID(1)
	want.SetBuildID(1)
	want.SetNumber(1)
	want.SetReporter("Foo Runtime")
	want.SetName("foo")
	want.SetMimetype("text/plain")

	got := new(library.InitStep)

	// setup database
	db, _ := sqlite.NewTest()

	defer func() {
		db.Sqlite.Exec("delete from repos;")
		db.Sqlite.Exec("delete from builds;")
		db.Sqlite.Exec("delete from initsteps;")
		_sql, _ := db.Sqlite.DB()
		_sql.Close()
	}()

	_ = db.CreateRepo(r)
	_ = db.CreateBuild(b)
	_ = db.CreateInitStep(want)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/1/initsteps/1", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(build.Establish())
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds/:build/initsteps/:initstep", func(c *gin.Context) {
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

func TestInitStep_Establish_NoRepo(t *testing.T) {
	// setup database
	db, _ := sqlite.NewTest()

	defer func() { _sql, _ := db.Sqlite.DB(); _sql.Close() }()

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/1/initsteps/1", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds/:build/initsteps/:initstep", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusNotFound)
	}
}

func TestInitStep_Establish_NoBuild(t *testing.T) {
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
	db, _ := sqlite.NewTest()

	defer func() {
		db.Sqlite.Exec("delete from repos;")
		_sql, _ := db.Sqlite.DB()
		_sql.Close()
	}()

	_ = db.CreateRepo(r)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/1/initsteps/1", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds/:build/initsteps/:initstep", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusNotFound)
	}
}

func TestInitStep_Establish_NoInitStepParameter(t *testing.T) {
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
	db, _ := sqlite.NewTest()

	defer func() {
		db.Sqlite.Exec("delete from repos;")
		db.Sqlite.Exec("delete from builds;")
		_sql, _ := db.Sqlite.DB()
		_sql.Close()
	}()

	_ = db.CreateRepo(r)
	_ = db.CreateBuild(b)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/1/initsteps", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(build.Establish())
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds/:build/initsteps", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusBadRequest)
	}
}

func TestInitStep_Establish_InvalidInitStepParameter(t *testing.T) {
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
	db, _ := sqlite.NewTest()

	defer func() {
		db.Sqlite.Exec("delete from repos;")
		db.Sqlite.Exec("delete from builds;")
		_sql, _ := db.Sqlite.DB()
		_sql.Close()
	}()

	_ = db.CreateRepo(r)
	_ = db.CreateBuild(b)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/1/initsteps/foo", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(build.Establish())
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds/:build/initsteps/:initstep", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusBadRequest)
	}
}

func TestInitStep_Establish_NoInitStep(t *testing.T) {
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
	db, _ := sqlite.NewTest()

	defer func() {
		db.Sqlite.Exec("delete from repos;")
		db.Sqlite.Exec("delete from builds;")
		_sql, _ := db.Sqlite.DB()
		_sql.Close()
	}()

	_ = db.CreateRepo(r)
	_ = db.CreateBuild(b)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/1/initsteps/1", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(build.Establish())
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds/:build/initsteps/:initstep", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusNotFound)
	}
}
