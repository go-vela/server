// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/go-vela/server/storage"
	"github.com/go-vela/server/storage/minio"
)

func TestMiddleware_Storage(t *testing.T) {
	// setup types
	var got storage.Storage

	want, _ := minio.NewTest("", "", "", "", false)
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(Storage(want))
	engine.GET("/health", func(c *gin.Context) {
		got = storage.FromGinContext(c)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Storage returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Storage is %v, want %v", got, want)
	}
}

func TestMiddleware_StorageEnable(t *testing.T) {
	// setup types
	got := false
	want := true
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodGet, "/health", nil)
	// setup mock server
	engine.Use(StorageEnable(want))
	engine.GET("/health", func(c *gin.Context) {
		got = c.Value("storage-enable").(bool)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("StorageEnable returned %v, want %v", resp.Code, http.StatusOK)
	}

	if got != want {
		t.Errorf("StorageEnable is %v, want %v", got, want)
	}
}
