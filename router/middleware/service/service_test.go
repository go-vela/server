// SPDX-License-Identifier: Apache-2.0

package service

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
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
)

func TestService_Retrieve(t *testing.T) {
	// setup types
	want := new(api.Service)
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

func TestService_Establish(t *testing.T) {
	// setup types
	owner := new(api.User)
	owner.SetID(1)

	r := new(api.Repo)
	r.SetID(1)
	r.SetOwner(owner)
	r.SetHash("baz")
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")
	r.SetVisibility("public")

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)

	want := new(api.Service)
	want.SetID(1)
	want.SetRepoID(1)
	want.SetBuildID(1)
	want.SetNumber(1)
	want.SetName("foo")
	want.SetImage("baz")
	want.SetStatus("")
	want.SetError("")
	want.SetExitCode(0)
	want.SetCreated(0)
	want.SetStarted(0)
	want.SetFinished(0)
	want.SetHost("")
	want.SetRuntime("")
	want.SetDistribution("")

	got := new(api.Service)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteBuild(context.TODO(), b)
		_ = db.DeleteRepo(context.TODO(), r)
		_ = db.DeleteService(context.TODO(), want)
		db.Close()
	}()

	_, _ = db.CreateRepo(context.TODO(), r)
	_, _ = db.CreateBuild(context.TODO(), b)
	_, _ = db.CreateService(context.TODO(), want)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/1/services/1", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(build.Establish())
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds/:build/services/:service", func(c *gin.Context) {
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

func TestService_Establish_NoRepo(t *testing.T) {
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
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/1/services/1", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds/:build/services/:service", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusNotFound)
	}
}

func TestService_Establish_NoBuild(t *testing.T) {
	// setup types
	r := new(api.Repo)
	r.SetID(1)
	r.GetOwner().SetID(1)
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
		_ = db.DeleteRepo(context.TODO(), r)
		db.Close()
	}()

	_, _ = db.CreateRepo(context.TODO(), r)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/1/services/1", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds/:build/services/:service", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusNotFound)
	}
}

func TestService_Establish_NoServiceParameter(t *testing.T) {
	// setup types
	owner := new(api.User)
	owner.SetID(1)

	r := new(api.Repo)
	r.SetID(1)
	r.SetOwner(owner)
	r.SetHash("baz")
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")
	r.SetVisibility("public")

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteBuild(context.TODO(), b)
		_ = db.DeleteRepo(context.TODO(), r)
		db.Close()
	}()

	_, _ = db.CreateRepo(context.TODO(), r)
	_, _ = db.CreateBuild(context.TODO(), b)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/1/services", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(build.Establish())
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds/:build/services", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusBadRequest)
	}
}

func TestService_Establish_InvalidServiceParameter(t *testing.T) {
	// setup types
	owner := new(api.User)
	owner.SetID(1)

	r := new(api.Repo)
	r.SetID(1)
	r.SetOwner(owner)
	r.SetHash("baz")
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")
	r.SetVisibility("public")

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteBuild(context.TODO(), b)
		_ = db.DeleteRepo(context.TODO(), r)
		db.Close()
	}()

	_, _ = db.CreateRepo(context.TODO(), r)
	_, _ = db.CreateBuild(context.TODO(), b)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/1/services/foo", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(build.Establish())
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds/:build/services/:service", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusBadRequest)
	}
}

func TestService_Establish_NoService(t *testing.T) {
	// setup types
	r := new(api.Repo)
	r.SetID(1)
	r.GetOwner().SetID(1)
	r.SetHash("baz")
	r.SetOrg("foo")
	r.SetName("bar")
	r.SetFullName("foo/bar")
	r.SetVisibility("public")

	b := new(api.Build)
	b.SetID(1)
	b.SetRepo(r)
	b.SetNumber(1)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteBuild(context.TODO(), b)
		_ = db.DeleteRepo(context.TODO(), r)
		db.Close()
	}()

	_, _ = db.CreateRepo(context.TODO(), r)
	_, _ = db.CreateBuild(context.TODO(), b)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/foo/bar/builds/1/services/1", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(org.Establish())
	engine.Use(repo.Establish())
	engine.Use(build.Establish())
	engine.Use(Establish())
	engine.GET("/:org/:repo/builds/:build/services/:service", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusNotFound)
	}
}
