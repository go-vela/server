// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package org

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/types/library"
)

func TestOrg_Retrieve(t *testing.T) {
	// setup types
	want := ""

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

func TestOrg_Establish(t *testing.T) {
	// setup types
	r := new(library.Repo)
	r.SetID(1)
	r.SetUserID(1)
	r.SetHash("baz")
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")
	r.SetLink("")
	r.SetClone("")
	r.SetBranch("")
	r.SetTimeout(0)
	r.SetVisibility("public")
	r.SetPrivate(false)
	r.SetTrusted(false)
	r.SetActive(false)
	r.SetAllowEvents(library.NewEventsFromMask(1))

	want := "foo"
	got := ""

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		db.DeleteRepo(r)
		db.Close()
	}()

	_ = db.CreateRepo(r)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(Establish())
	engine.GET("/:org", func(c *gin.Context) {
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

func TestOrg_Establish_NoOrgParameter(t *testing.T) {
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
	context.Request, _ = http.NewRequest(http.MethodGet, "//test", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(Establish())
	engine.GET("/:org/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusBadRequest)
	}
}
