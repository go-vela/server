// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package middleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMiddleware_QueueSigningPrivateKey(t *testing.T) {
	// setup types
	got := ""
	want := "foobar"

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(QueueSigningPrivateKey(want))
	engine.GET("/health", func(c *gin.Context) {
		got = c.Value("queue.private-key").(string)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("QueueSigningPrivateKey returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("QueueSigningPrivateKey is %v, want %v", got, want)
	}
}

func TestMiddleware_QueueSigningPublicKey(t *testing.T) {
	// setup types
	got := ""
	want := "foobar"

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(QueueSigningPublicKey(want))
	engine.GET("/health", func(c *gin.Context) {
		got = c.Value("public-key").(string)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("QueueSigningPublicKey returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("QueueSigningPublicKey is %v, want %v", got, want)
	}
}

func TestMiddleware_QueueAddress(t *testing.T) {
	// setup types
	got := ""
	want := "foobar"

	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	context, engine := gin.CreateTestContext(resp)
	context.Request, _ = http.NewRequest(http.MethodGet, "/health", nil)

	// setup mock server
	engine.Use(QueueAddress(want))
	engine.GET("/health", func(c *gin.Context) {
		got = c.Value("queue-address").(string)

		c.Status(http.StatusOK)
	})

	// run test
	engine.ServeHTTP(context.Writer, context.Request)

	if resp.Code != http.StatusOK {
		t.Errorf("QueueAddress returned %v, want %v", resp.Code, http.StatusOK)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("QueueAddress is %v, want %v", got, want)
	}
}
