// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/types/library"
)

func TestDashboard_FromContext(t *testing.T) {
	// setup types
	uuid := "c8da1302-07d6-11ea-882f-4893bca275b8"
	want := &library.Dashboard{ID: &uuid}

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

func TestDashboard_FromContext_Bad(t *testing.T) {
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

func TestDashboard_FromContext_WrongType(t *testing.T) {
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

func TestDashboard_FromContext_Empty(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(nil)

	// run test
	got := FromContext(context)

	if got != nil {
		t.Errorf("FromContext is %v, want nil", got)
	}
}

func TestDashboard_ToContext(t *testing.T) {
	// setup types
	uuid := "c8da1302-07d6-11ea-882f-4893bca275b8"
	want := &library.Dashboard{ID: &uuid}

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
