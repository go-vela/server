// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v3"
)

func TestMiddleware_CLI(t *testing.T) {
	// setup types
	want := &cli.Context{
		App: &cli.App{
			Name: "foo",
		},
	}

	got := &cli.Context{}

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(CLI(want))
	engine.GET("/health", func(c *gin.Context) {
		got = c.Value("cli").(*cli.Context)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("CLI returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("CLI is %v, want %v", got, want)
	}
}
