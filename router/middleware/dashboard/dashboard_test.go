// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/types/library"
)

func TestDashboard_Retrieve(t *testing.T) {
	// setup types
	want := new(library.Dashboard)
	want.SetID("c8da1302-07d6-11ea-882f-4893bca275b8")

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

func TestDashboard_Establish(t *testing.T) {
	// setup types
	want := new(library.Dashboard)
	want.SetID("c8da1302-07d6-11ea-882f-4893bca275b8")
	want.SetName("vela")
	want.SetCreatedAt(1)
	want.SetCreatedBy("octocat")
	want.SetUpdatedAt(1)
	want.SetUpdatedBy("octokitty")
	want.SetAdmins([]string{"octocat", "octokitty"})

	wantRepo := new(library.DashboardRepo)
	wantRepo.SetID(1)
	wantRepo.SetName("go-vela/server")
	wantRepo.SetBranches([]string{"main"})
	wantRepo.SetEvents([]string{"push"})

	want.SetRepos([]*library.DashboardRepo{wantRepo})

	got := new(library.Dashboard)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		_ = db.DeleteDashboard(context.TODO(), want)
		db.Close()
	}()

	_, _ = db.CreateDashboard(context.TODO(), want)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/c8da1302-07d6-11ea-882f-4893bca275b8", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(Establish())
	engine.GET("/:dashboard", func(c *gin.Context) {
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

func TestDashboard_Establish_NoDashboardParameter(t *testing.T) {
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
	engine.GET("/:dashboard/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusBadRequest)
	}
}

func TestDashboard_Establish_NoDashboard(t *testing.T) {
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
	context.Request, _ = http.NewRequest(http.MethodGet, "/c8da1302-07d6-11ea-882f-4893bca275b8", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(Establish())
	engine.GET("/:dashboard", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusNotFound)
	}
}
