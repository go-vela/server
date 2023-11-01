// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/types/library"
)

func TestWorker_Retrieve(t *testing.T) {
	// setup types
	want := new(api.Worker)
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

func TestWorker_Establish(t *testing.T) {
	// setup types
	b := new(library.Build)
	b.SetID(1)

	want := new(api.Worker)
	want.SetID(1)
	want.SetHostname("worker_0")
	want.SetAddress("localhost")
	want.SetRoutes([]string{"foo", "bar", "baz"})
	want.SetActive(true)
	want.SetStatus("available")
	want.SetLastStatusUpdateAt(12345)
	want.SetRunningBuilds([]*library.Build{b})
	want.SetLastBuildStartedAt(12345)
	want.SetLastBuildFinishedAt(12345)
	want.SetLastCheckedIn(12345)
	want.SetBuildLimit(0)

	got := new(api.Worker)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	defer func() {
		db.DeleteWorker(context.TODO(), want)
		db.Close()
	}()

	_, _ = db.CreateWorker(context.TODO(), want)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/workers/worker_0", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(Establish())
	engine.GET("/workers/:worker", func(c *gin.Context) {
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

func TestWorker_Establish_NoWorkerParameter(t *testing.T) {
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
	context.Request, _ = http.NewRequest(http.MethodGet, "/workers/", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(Establish())
	engine.GET("/workers/:worker", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Establish returned %v, want %v", resp.Code, http.StatusBadRequest)
	}
}
