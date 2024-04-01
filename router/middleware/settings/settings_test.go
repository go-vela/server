// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/types/library"
)

func TestSettings_Retrieve(t *testing.T) {
	// setup types
	want := new(api.Settings)
	// want.SetID(1)

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

func TestSettings_Establish(t *testing.T) {
	// setup types
	b := new(library.Build)
	b.SetID(1)

	want := new(api.Settings)
	// want.SetID(1)
	// want.SetHostname("worker_0")
	// want.SetAddress("localhost")
	// want.SetRoutes([]string{"foo", "bar", "baz"})
	// want.SetActive(true)
	// want.SetStatus("available")
	// want.SetLastStatusUpdateAt(12345)
	// want.SetRunningBuilds([]*library.Build{b})
	// want.SetLastBuildStartedAt(12345)
	// want.SetLastBuildFinishedAt(12345)
	// want.SetLastCheckedIn(12345)
	// want.SetBuildLimit(0)

	got := new(api.Settings)

	// setup database
	db, err := database.NewTest()
	if err != nil {
		t.Errorf("unable to create test database engine: %v", err)
	}

	// defer func() {
	// 	_ = db.DeleteSettings(context.TODO(), want)
	// 	db.Close()
	// }()

	// _, _ = db.CreateSettings(context.TODO(), want)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/admin/settings", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { database.ToContext(c, db) })
	engine.Use(Establish())
	engine.GET("/admin/settings", func(c *gin.Context) {
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
