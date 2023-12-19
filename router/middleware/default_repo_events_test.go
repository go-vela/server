// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-vela/types/constants"

	"github.com/gin-gonic/gin"
)

func TestMiddleware_DefaultRepoEvents(t *testing.T) {
	// setup types
	var got []string

	want := []string{constants.EventPush}

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(DefaultRepoEvents(want))
	engine.GET("/health", func(c *gin.Context) {
		got = c.Value("defaultRepoEvents").([]string)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("DefaultRepoEvents returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("DefaultRepoEvents is %v, want %v", got, want)
	}
}

func TestMiddleware_DefaultRepoEventsMask(t *testing.T) {
	// setup types
	var got int64

	want := int64(3)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(DefaultRepoEventsMask(want))
	engine.GET("/health", func(c *gin.Context) {
		got = c.Value("defaultRepoEventsMask").(int64)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("DefaultRepoEventsMask returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("DefaultRepoEventsMask is %v, want %v", got, want)
	}
}
