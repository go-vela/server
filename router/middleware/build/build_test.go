// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/repo"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
)

func TestBuild_Retrieve(t *testing.T) {
	// setup types
	bID := int64(1)
	want := &library.Build{ID: &bID}

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

func TestBuild_Establish(t *testing.T) {
	// setup types
	uID := int64(1)
	uUserID := int64(1)
	uOrg := "foo"
	uName := "bar"
	uFullName := "foo/bar"
	r := &library.Repo{
		ID:       &uID,
		UserID:   &uUserID,
		Org:      &uOrg,
		Name:     &uName,
		FullName: &uFullName,
	}

	bID := int64(1)
	bRepoID := int64(1)
	bNumber := 1
	zeroInt := 0
	zeroInt64 := int64(0)
	zeroString := ""
	want := &library.Build{
		ID:           &bID,
		RepoID:       &bRepoID,
		Number:       &bNumber,
		Parent:       &zeroInt,
		Event:        &zeroString,
		Status:       &zeroString,
		Error:        &zeroString,
		Enqueued:     &zeroInt64,
		Created:      &zeroInt64,
		Started:      &zeroInt64,
		Finished:     &zeroInt64,
		Deploy:       &zeroString,
		Clone:        &zeroString,
		Source:       &zeroString,
		Title:        &zeroString,
		Message:      &zeroString,
		Commit:       &zeroString,
		Sender:       &zeroString,
		Author:       &zeroString,
		Branch:       &zeroString,
		Ref:          &zeroString,
		BaseRef:      &zeroString,
		Host:         &zeroString,
		Runtime:      &zeroString,
		Distribution: &zeroString,
	}
	got := new(library.Build)

	// setup database
	db, _ := database.NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Exec("delete from builds;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(r)
	_ = db.CreateBuild(want)

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/1", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(repo.Establish())
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds/:build", func(c *gin.Context) {
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

func TestBuild_Establish_NoRepo(t *testing.T) {
	// setup database
	db, _ := database.NewTest()
	defer db.Database.Close()

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/1", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(Establish())

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusNotFound)
	}
}

func TestBuild_Establish_NoBuildParameter(t *testing.T) {
	// setup types
	rID := int64(1)
	rUserID := int64(1)
	rOrg := "foo"
	rName := "bar"
	rFullName := "foo/bar"
	r := &library.Repo{
		ID:       &rID,
		UserID:   &rUserID,
		Org:      &rOrg,
		Name:     &rName,
		FullName: &rFullName,
	}

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
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(repo.Establish())
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusBadRequest)
	}
}

func TestBuild_Establish_InvalidBuildParameter(t *testing.T) {
	// setup types
	rID := int64(1)
	rUserID := int64(1)
	rOrg := "foo"
	rName := "bar"
	rFullName := "foo/bar"
	r := &library.Repo{
		ID:       &rID,
		UserID:   &rUserID,
		Org:      &rOrg,
		Name:     &rName,
		FullName: &rFullName,
	}

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
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/foo", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(repo.Establish())
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds/:build", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusBadRequest)
	}
}

func TestBuild_Establish_NoBuild(t *testing.T) {
	// setup types
	rID := int64(1)
	rUserID := int64(1)
	rOrg := "foo"
	rName := "bar"
	rFullName := "foo/bar"
	r := &library.Repo{
		ID:       &rID,
		UserID:   &rUserID,
		Org:      &rOrg,
		Name:     &rName,
		FullName: &rFullName,
	}

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
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/1", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(repo.Establish())
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds/:build", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusNotFound)
	}
}
