// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-vela/server/internal/token"

	"github.com/gin-gonic/gin"
)

func TestMiddleware_TokenManager(t *testing.T) {
	// setup types
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	var got *token.Manager

	want := new(token.Manager)
	want.PrivateKey = "123abc"

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(TokenManager(want))
	engine.GET("/health", func(c *gin.Context) {
		got = c.MustGet("token-manager").(*token.Manager)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("TokenManager returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("TokenManager is %v, want %v", got, want)
	}
}
