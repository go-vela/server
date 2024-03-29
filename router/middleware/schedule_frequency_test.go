// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestMiddleware_ScheduleFrequency(t *testing.T) {
	// setup types
	var got time.Duration
	want := 30 * time.Minute

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(ScheduleFrequency(want))
	engine.GET("/health", func(c *gin.Context) {
		got = c.Value("scheduleminimumfrequency").(time.Duration)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("ScheduleFrequency returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ScheduleFrequency is %v, want %v", got, want)
	}
}
