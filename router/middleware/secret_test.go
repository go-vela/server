// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-vela/server/database"

	"github.com/go-vela/server/secret"
	"github.com/go-vela/server/secret/native"

	"github.com/gin-gonic/gin"
)

func TestMiddleware_Secret(t *testing.T) {
	// setup types
	got := ""
	want := "foobar"

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(Secret(want))
	engine.GET("/health", func(c *gin.Context) {
		got = c.Value("secret").(string)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Secret returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Secret is %v, want %v", got, want)
	}
}

func TestMiddleware_Secrets(t *testing.T) {
	// setup types
	d, _ := database.NewTest()
	defer d.Database.Close()

	var got secret.Service

	want, _ := native.New(d)
	s := map[string]secret.Service{"native": want}

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(Secrets(s))
	engine.GET("/health", func(c *gin.Context) {
		got = secret.FromContext(c, "native")

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("Secrets returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Secrets is %v, want %v", got, want)
	}
}
