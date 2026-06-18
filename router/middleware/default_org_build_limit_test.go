// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMiddleware_DefaultOrgBuildLimit(t *testing.T) {
	// setup types
	var got int32

	want := int32(30)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(DefaultOrgBuildLimit(want))
	engine.GET("/health", func(c *gin.Context) {
		got = c.Value("defaultOrgBuildLimit").(int32)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("DefaultOrgBuildLimit returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("DefaultOrgBuildLimit is %v, want %v", got, want)
	}
}
