// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
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
	num := int64(1)
	want := &library.Step{ID: &num}

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
	num := 1
	num64 := int64(1)
	foo := "foo"
	bar := "bar"
	foobar := "foo/bar"
	zeroInt := 1
	zeroInt64 := int64(1)
	zeroString := ""
	r := &library.Repo{ID: &num64, UserID: &num64, Org: &foo, Name: &bar, FullName: &foobar}
	b := &library.Build{ID: &num64, RepoID: &num64, Number: &num}
	want := &library.Step{
		ID:           &num64,
		RepoID:       &num64,
		BuildID:      &num64,
		Number:       &num,
		Name:         &foo,
		Stage:        &zeroString,
		Status:       &zeroString,
		Error:        &zeroString,
		ExitCode:     &zeroInt,
		Created:      &zeroInt64,
		Started:      &zeroInt64,
		Finished:     &zeroInt64,
		Host:         &zeroString,
		Runtime:      &zeroString,
		Distribution: &zeroString,
	}
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
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
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
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
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
	rID := int64(1)
	rOrg := "foo"
	rName := "bar"
	rFullName := "foo/bar"
	r := &library.Repo{ID: &rID, UserID: &rID, Org: &rOrg, Name: &rName, FullName: &rFullName}

	// setup database
	db, _ := database.NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(r)

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
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
	num := 1
	num64 := int64(1)
	rOrg := "foo"
	rName := "bar"
	rFullName := "foo/bar"
	r := &library.Repo{ID: &num64, UserID: &num64, Org: &rOrg, Name: &rName, FullName: &rFullName}
	b := &library.Build{ID: &num64, RepoID: &num64, Number: &num}

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
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
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
	num := 1
	num64 := int64(1)
	rOrg := "foo"
	rName := "bar"
	rFullName := "foo/bar"
	r := &library.Repo{ID: &num64, UserID: &num64, Org: &rOrg, Name: &rName, FullName: &rFullName}
	b := &library.Build{ID: &num64, RepoID: &num64, Number: &num}

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
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
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
	num := 1
	num64 := int64(1)
	rOrg := "foo"
	rName := "bar"
	rFullName := "foo/bar"
	r := &library.Repo{ID: &num64, UserID: &num64, Org: &rOrg, Name: &rName, FullName: &rFullName}
	b := &library.Build{ID: &num64, RepoID: &num64, Number: &num}

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
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
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
