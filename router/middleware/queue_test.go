// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/queue/redis"

	"github.com/gin-gonic/gin"
)

func TestMiddleware_Queue(t *testing.T) {
	// setup types
	var got queue.Service

	// signing keys are irrelevant here
	want, _ := redis.NewTest("", "")

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(Queue(want))
	engine.GET("/health", func(c *gin.Context) {
		got = queue.FromGinContext(c)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Queue returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Queue is %v, want %v", got, want)
	}
}
