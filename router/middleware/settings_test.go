// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/api/types/settings"
)

func TestMiddleware_Settings(t *testing.T) {
	// setup types
	want := settings.PlatformMockEmpty()
	want.SetCloneImage("target/vela-git")

	got := settings.PlatformMockEmpty()

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(Settings(&want))
	engine.GET("/health", func(c *gin.Context) {
		got = *c.Value("settings").(*settings.Platform)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Settings returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Settings is %v, want %v", got, want)
	}
}
