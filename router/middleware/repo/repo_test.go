// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-vela/server/router/middleware/org"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/database/sqlite"
	"github.com/go-vela/types/library"
)

func TestRepo_Retrieve(t *testing.T) {
	// setup types
	want := new(library.Repo)
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

func TestRepo_Establish(t *testing.T) {
	// setup types
	want := new(library.Repo)
	want.SetID(1)
	want.SetUserID(1)
	want.SetHash("baz")
	want.SetOrg("foo")
	want.SetName("bar")
	want.SetFullName("foo/bar")
	want.SetLink("")
	want.SetClone("")
	want.SetBranch("")
	want.SetBuildLimit(0)
	want.SetTimeout(0)
	want.SetCounter(0)
	want.SetVisibility("public")
	want.SetPrivate(false)
	want.SetTrusted(false)
	want.SetActive(false)
	want.SetAllowPull(false)
	want.SetAllowPush(false)
	want.SetAllowDeploy(false)
	want.SetAllowTag(false)
	want.SetAllowComment(false)
	want.SetPipelineType("yaml")
	want.SetPreviousName("")
	want.SetLastUpdate(0)

	got := new(library.Repo)

	// setup database
	db, _ := sqlite.NewTest()

	defer func() {
		db.Sqlite.Exec("delete from repos;")
		_sql, _ := db.Sqlite.DB()
		_sql.Close()
	}()

	_ = db.CreateRepo(want)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(org.Establish())
	engine.Use(Establish())
	engine.GET("/:org/:repo", func(c *gin.Context) {
		got = Retrieve(c)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(resp, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Establish is %v, want %v", got, want)
	}
}

func TestRepo_Establish_NoOrgParameter(t *testing.T) {
	// setup database
	db, _ := sqlite.NewTest()

	defer func() { _sql, _ := db.Sqlite.DB(); _sql.Close() }()

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "//bar/test", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(Establish())
	engine.GET("/:org/:repo/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusBadRequest)
	}
}

func TestRepo_Establish_NoRepoParameter(t *testing.T) {
	// setup database
	db, _ := sqlite.NewTest()

	defer func() { _sql, _ := db.Sqlite.DB(); _sql.Close() }()

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo//test", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(Establish())
	engine.GET("/:org/:repo/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusBadRequest)
	}
}

func TestRepo_Establish_NoRepo(t *testing.T) {
	// setup database
	db, _ := sqlite.NewTest()

	defer func() { _sql, _ := db.Sqlite.DB(); _sql.Close() }()

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(Establish())
	engine.GET("/:org/:repo", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusNotFound)
	}
}
