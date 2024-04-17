// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"

	"github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/compiler/native"
)

func TestMiddleware_CompilerNative(t *testing.T) {
	// setup types
	var got compiler.Engine

	wantCloneImage := "target/vela-git"
	want, _ := native.New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))
	want.SetCloneImage(wantCloneImage)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(Compiler(want))
	// setup mock server
	engine.Use(func(c *gin.Context) {
		s := new(settings.Platform)
		// todo: this should fail
		// s.SetCloneImage(wantCloneImage)

		c.Set("settings", s)
		c.Next()
	})
	engine.GET("/health", func(c *gin.Context) {
		got = compiler.FromContext(c)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Compiler returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Compiler is %v, want %v", got, want)
	}
}
