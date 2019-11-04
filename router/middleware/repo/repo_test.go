// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-vela/server/database"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
)

func TestRepo_Retrieve(t *testing.T) {
	// setup types
	num := int64(1)
	want := &library.Repo{ID: &num}

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
	num := int64(1)
	org := "foo"
	name := "bar"
	fullname := "foo/bar"
	zeroInt64 := int64(0)
	zeroString := ""
	zeroBool := false
	want := &library.Repo{
		ID:          &num,
		UserID:      &num,
		Org:         &org,
		Name:        &name,
		FullName:    &fullname,
		Link:        &zeroString,
		Clone:       &zeroString,
		Branch:      &zeroString,
		Timeout:     &zeroInt64,
		Visibility:  &zeroString,
		Private:     &zeroBool,
		Trusted:     &zeroBool,
		Active:      &zeroBool,
		AllowPull:   &zeroBool,
		AllowPush:   &zeroBool,
		AllowDeploy: &zeroBool,
		AllowTag:    &zeroBool,
	}
	got := new(library.Repo)

	// setup database
	db, _ := database.NewTest()
	defer func() {
		db.Database.Exec("delete from repos;")
		db.Database.Close()
	}()
	_ = db.CreateRepo(want)

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
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
	db, _ := database.NewTest()
	defer db.Database.Close()

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
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
	db, _ := database.NewTest()
	defer db.Database.Close()

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
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
	db, _ := database.NewTest()
	defer db.Database.Close()

	// setup context
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
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
