// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/compiler/native"

	"github.com/gin-gonic/gin"

	"github.com/urfave/cli/v2"
)

func TestMiddleware_CompilerNative(t *testing.T) {
	// setup types
	var got compiler.Engine

	want, _ := native.New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(Compiler(want))
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
