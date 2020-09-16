// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package source

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-vela/server/source/github"

	"github.com/gin-gonic/gin"
)

func TestSource_FromContext(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	context, engine := gin.CreateTestContext(nil)

	// setup mock server
	engine.GET("/health", func(c *gin.Context) { c.String(http.StatusOK, "OK-") })

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	want, _ := github.NewTest(s.URL)
	context.Set(key, want)

	// run test
	got := FromContext(context)

	if got != want {
		t.Errorf("FromContext is %v, want %v", got, want)
	}
}

func TestSource_FromContext_Bad(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	context, _ := gin.CreateTestContext(nil)
	context.Set(key, nil)

	// run test
	got := FromContext(context)

	if got != nil {
		t.Errorf("FromContext is %v, want nil", got)
	}
}

func TestSource_FromContext_WrongType(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	context, _ := gin.CreateTestContext(nil)
	context.Set(key, 1)

	// run test
	got := FromContext(context)

	if got != nil {
		t.Errorf("FromContext is %v, want nil", got)
	}
}

func TestSource_FromContext_Empty(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	context, _ := gin.CreateTestContext(nil)

	// run test
	got := FromContext(context)

	if got != nil {
		t.Errorf("FromContext is %v, want nil", got)
	}
}

func TestSource_ToContext(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	context, engine := gin.CreateTestContext(nil)

	// setup mock server
	engine.GET("/health", func(c *gin.Context) { c.String(http.StatusOK, "OK") })

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	want, _ := github.NewTest(s.URL)
	ToContext(context, want)

	// run test
	got := context.Value(key)

	if got != want {
		t.Errorf("ToContext is %v, want %v", got, want)
	}
}
