// SPDX-License-Identifier: Apache-2.0

package org

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestOrg_FromContext(t *testing.T) {
	// setup types
	want := "foo"

	// setup context
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(nil)
	context.Set(key, want)

	// run test
	got := FromContext(context)

	if got != want {
		t.Errorf("FromContext is %v, want %v", got, want)
	}
}

func TestOrg_FromContext_Bad(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(nil)
	context.Set(key, nil)

	// run test
	got := FromContext(context)

	if got != "" {
		t.Errorf("FromContext is %v, want nil", got)
	}
}

func TestOrg_FromContext_WrongType(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(nil)
	context.Set(key, 1)

	// run test
	got := FromContext(context)

	if got != "" {
		t.Errorf("FromContext is %v, want nil", got)
	}
}

func TestOrg_FromContext_Empty(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(nil)

	// run test
	got := FromContext(context)

	if got != "" {
		t.Errorf("FromContext is %v, want nil", got)
	}
}

func TestOrg_ToContext(t *testing.T) {
	// setup types
	want := "foo"

	// setup context
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(nil)
	ToContext(context, want)

	// run test
	got := context.Value(key)

	if got != want {
		t.Errorf("ToContext is %v, want %v", got, want)
	}
}
