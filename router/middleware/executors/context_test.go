// SPDX-License-Identifier: Apache-2.0

package executors

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	api "github.com/go-vela/server/api/types"
)

func TestExecutors_FromContext(t *testing.T) {
	// setup types
	eID := int64(1)
	e := api.Executor{ID: &eID}
	want := []api.Executor{e}

	// setup context
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(nil)
	context.Set(key, want)

	// run test
	got := FromContext(context)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("FromContext is %v, want %v", got, want)
	}
}

func TestExecutors_FromContext_Bad(t *testing.T) {
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

func TestExecutors_FromContext_WrongType(t *testing.T) {
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

func TestExecutors_FromContext_Empty(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(nil)

	// run test
	got := FromContext(context)

	if got != nil {
		t.Errorf("FromContext is %v, want nil", got)
	}
}

func TestExecutors_ToContext(t *testing.T) {
	// setup types
	eID := int64(1)
	e := api.Executor{ID: &eID}
	want := []api.Executor{e}

	// setup context
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(nil)
	ToContext(context, want)

	// run test
	got := context.Value(key)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToContext is %v, want %v", got, want)
	}
}
