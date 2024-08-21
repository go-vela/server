// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/tracing"
)

func TestTracing_Retrieve(t *testing.T) {
	// setup types
	serviceName := "vela-test"
	want := &tracing.Client{
		Config: tracing.Config{
			ServiceName: serviceName,
		},
	}

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

func TestTracing_Establish(t *testing.T) {
	// setup types
	serviceName := "vela-test"
	want := &tracing.Client{
		Config: tracing.Config{
			ServiceName: serviceName,
		},
	}

	got := new(tracing.Client)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/hello", nil)

	// setup mock server
	engine.Use(func(c *gin.Context) { c.Set("logger", logrus.NewEntry(logrus.StandardLogger())) })
	engine.Use(func(c *gin.Context) { ToContext(c, want) })
	engine.Use(Establish())
	engine.GET("/hello", func(c *gin.Context) {
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
