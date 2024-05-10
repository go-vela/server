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
	sMiddleware "github.com/go-vela/server/router/middleware/settings"
)

func TestMiddleware_CompilerNative(t *testing.T) {
	// setup types
	defaultCloneImage := "target/vela-git"
	wantCloneImage := "target/vela-git:latest"

	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.String("clone-image", "", "")
	fs.Set("clone-image", defaultCloneImage)

	cliCtx := cli.NewContext(nil, fs, nil)

	want, _ := native.FromCLIContext(cliCtx)
	want.SetCloneImage(wantCloneImage)

	var got compiler.Engine
	got, _ = native.FromCLIContext(cliCtx)

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)

	engine.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			s := settings.Platform{
				Compiler: &settings.Compiler{},
			}
			s.SetCloneImage(wantCloneImage)

			sMiddleware.ToContext(c, &s)

			c.Next()
		}
	}(),
	)

	engine.Use(Compiler(got))

	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

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
