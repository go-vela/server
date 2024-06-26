// SPDX-License-Identifier: Apache-2.0

package org

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
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
	r := new(api.Repo)
	r.SetID(1)
	r.GetOwner().SetID(1)
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

	want := "foo"
	got := ""

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteRepo(context.TODO(), r)
		db.Close()
	}()

	_, _ = db.CreateRepo(context.TODO(), r)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
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
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
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
